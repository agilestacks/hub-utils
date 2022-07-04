// Copyright (c) 2022 EPAM Systems, Inc.
// 
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

const functions = require('@google-cloud/functions-framework');
const { Storage } = require('@google-cloud/storage');
const { ungzip } = require('node-gzip');
const YAML = require('js-yaml');
const request = require('request');

const storage = new Storage()

exports.reporting = async (file, context) => {
  if (context.eventType !== 'google.storage.object.finalize') {
    return
  }
  if (!file.name.includes('hub.state')) {
    return
  }
  const superhubs = await superhubBuckets()
  const stack = await stackByID(file.name, superhubs, false)
  let reportingEndpoint = 'https://us-central1-superhub.cloudfunctions.net/event'
  if (process.env.REPORTING_URL) {
    reportingEndpoint = process.env.REPORTING_URL
  }
  const payload = {
    stackId: stack.id,
    name: stack.name,
    status: stack.status,
    initiator: stack.latestOperation ? stack.latestOperation.initiator : 'unknown',
    project: stack.projectId,
    gcpUserAccount: stack.userAccount
  }
  request({
    url: reportingEndpoint,
    method: 'POST',
    json: payload
  }, function(error, response, body){
    if (response.statusCode = 200) {
      console.log(`Payload ${JSON.stringify(payload)} posted to ${reportingEndpoint}` )
    } else if (error) {
      console.log(error)
    }
  });
};

exports.stacks = async (req, res) => {
  const superhubs = await superhubBuckets()
  const id = stackID(req.path)
  switch (req.method) {
    case 'GET':
      if (id) {
        let stack
        if (req.query.raw || req.query.raw === '') {
          stack = await stackByID(id, superhubs, true)
        } else {
          stack = await stackByID(id, superhubs, false)
        }
        if (stack) {
          res
            .status(200)
            .set('content-type', 'application/json')
            .send(JSON.stringify(stack))
        } else {
          res
            .status(404)
            .send("Sorry, cant find that")
        }
      } else {
        const stacks = await allStacks(superhubs)
        res
          .status(200)
          .set('content-type', 'application/json')
          .send(JSON.stringify(filter(stacks,req.query)))
      }
      break
    case 'DELETE':
      stack = await stackByID(id, superhubs, false)
      if (stack) {
        await deleteByID(id, superhubs)
        res
          .status(202)
          .send("")
      } else {
        res
          .status(404)
          .send("")
      }
      break
    default:
      res
        .status(400)
        .send("Sorry, not supported")
  }
};

function filter(stacks, query) {
  let filtered = query.name ? stacks.filter(stack => stack.name.toLowerCase().includes(query.name.toLowerCase())) : stacks
  filtered = query['status'] ?
    filtered.filter(stack => stack.status.toLowerCase() === query['status'].toLowerCase()) : filtered
  filtered = query['latestOperation.name'] ?
    filtered.filter(stack => stack.latestOperation.name.toLowerCase() === query['latestOperation.name'].toLowerCase()) : filtered
  filtered = query['latestOperation.initiator'] ?
    filtered.filter(stack => stack.latestOperation.initiator.toLowerCase().includes(query['latestOperation.initiator'].toLowerCase())) : filtered
  if (query['latestOperation.timestamp']) {
    const before = query['latestOperation.timestamp'].before
    if (before) {
      const dtBefore = Date.parse(before)
      if (dtBefore) {
        filtered = filtered.filter(stack => new Date(stack.latestOperation.timestamp) < dtBefore)
      }
    }
  }
  if (query['latestOperation.timestamp']) {
    const after = query['latestOperation.timestamp'].after
    if (after) {
      const dtAfter = Date.parse(after)
      if (dtAfter) {
        filtered = filtered.filter(stack => new Date(stack.latestOperation.timestamp) >= dtAfter)
      }
    }
  }
  return filtered
}

async function deleteByID(id, buckets) {
  let stateFileMetas = []
  for (const bucket of buckets) {
    const [files] = await bucket.getFiles()
    stateFileMetas.push(...files.filter(bucketFile => bucketFile.name.includes(id)))
  }
  let promises = []
  for (const stateFileMeta of stateFileMetas) {
    promises.push(stateFileMeta.delete())
  }
  await Promise.all(promises)
}

async function stackByID(id, buckets, raw) {
  let stateFileMeta
  for (const bucket of buckets) {
    const [files] = await bucket.getFiles()
    const found = files
      .filter(bucketFile => bucketFile.name.includes(id) && bucketFile.name.includes('hub.state'))
    if (found.length > 0) {
      stateFileMeta = found[0]
      break
    }
  }
  if (!stateFileMeta) {
    return null
  }
  const states = await allStates([stateFileMeta])
  if (raw) {
    return YAML.load(states[0].toString())
  }
  return stackMeta(YAML.load(states[0].toString()), false)
}

async function allStacks(buckets) {
  let stateFileMetas = []
  for (const bucket of buckets) {
    const [files] = await bucket.getFiles()
    stateFileMetas.push(...files.filter(bucketFile => bucketFile.name.includes('hub.state')))
  }
  console.log('Number of state files: '+stateFileMetas.length)
  const start = new Date().getTime()
  const states = await allStates(stateFileMetas)
  const end = new Date().getTime()
  const time = (end - start)/1000
  console.log('Time spent to unzip all state files: '+time+' seconds')
  console.log('Average: '+time/states.length+' seconds per file')
  stacks = []
  states.forEach(state => {
    stacks.push(stackMeta(YAML.load(state), true))
  })
  stacks.sort((x,y) => new Date(y.latestOperation.timestamp) - new Date(x.latestOperation.timestamp))
  return stacks
}

function stackID(path) {
  if (!path || path === '/') {
    return ""
  } else if (path.startsWith('/')) {
    return path.substring(1)
  }
  return path
}

function stackMeta(yaml, light) {
  const last = yaml.operations[yaml.operations.length-1]
  let stateLocation = 'unknown'
  if (last.options.args) {
    const stateArg = last.options.args.find(option => option.includes('hub.state'))
    const locations = stateArg.split(',')
    if (locations.length === 2) {
      stateLocation = locations[1]
    }
  }
  let dnsDomainParam = {
    value: 'unset'
  }
  let projectId = {
    value: 'unset'
  }
  let userAccount = {
    value: last.initiator || 'unset'
  }
  let sandboxDir = {
    value: 'unset'
  }
  let sandboxCommit = {
    value: 'unset'
  }
  const {stackParameters} = yaml;
  if (stackParameters) {
    dnsDomainParam = stackParameters.find(({name}) => name === 'dns.domain') ||
      stackParameters.find(({name}) => name === 'dns.name') ||
      dnsDomainParam
    if (dnsDomainParam.value === 'unset' && stateLocation !== 'unknown') {
      const parts = stateLocation.split('/');
      const index = parts.indexOf('hub');
      if (index > 0) {
        dnsDomainParam.value = parts[index - 1];
      } else {
        dnsDomainParam.value = stateLocation;
      }
    }
    projectId = stackParameters.find(({name}) => name === 'projectId') || projectId
    userAccount = stackParameters.find(({name}) => name === 'hub.userAccount') || userAccount
    sandboxDir = stackParameters.find(({name}) => name === 'hub.sandboxDir') || sandboxDir
    sandboxCommit = stackParameters.find(({name}) => name === 'hub.sandboxCommit') || sandboxCommit
  }
  const meta = {
    id: dnsDomainParam.value,
    projectId: projectId.value,
    userAccount: userAccount.value,
    name: yaml.meta.name,
    stateLocation: {
      uri: stateLocation,
      kind: 'gcs'
    },
    sandbox: {
      dir: sandboxDir.value,
      commit: sandboxCommit.value
    },
    status: yaml.status,
    components: light ? undefined: Object.keys(yaml.components).map(key => {
      return {
          name: key,
          status: yaml.components[key].status,
          timestamp: yaml.components[key].timestamp
        }
      }
    ),
    latestOperation: {
      name: last.operation,
      timestamp: last.timestamp,
      status: last.status,
      initiator: userAccount.value, // by default value is last.initiator
      phases: last.phases
    }
  }
  return meta
}


async function allStates(stateFileMetas) {
  let promises = []
  for (const stateFileMeta of stateFileMetas) {
    promises.push(stateFile(stateFileMeta).then(archive => ungzip(archive)))
  }
  return Promise.all(promises)
}

function stateFile(file) {
  return new Promise((resolve, reject) => {
    let data = []
    file.createReadStream()
      .on('data', d => {
        data.push(d)
      })
      .on('end', () => {
        resolve(Buffer.concat(data));
      })
      .on('error', e => reject(e))
  })
}

async function superhubBuckets() {
  const [buckets] = await storage.getBuckets()
  if (!buckets || buckets.length === 0) {
    res.send('Project does not have any buckets!')
  }
  return buckets.filter(bucket => {
    const labels = bucket.metadata.labels
    if (labels) {
      const manager = labels['managed-by']
      return manager && manager === 'superhub'
    }
    return false
  })
}
