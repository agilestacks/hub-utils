const Knex = require('knex');
const {SecretManagerServiceClient} = require('@google-cloud/secret-manager');

const client = new SecretManagerServiceClient();

async function accessSecretVersion(secretName) {
  const [version] = await client.accessSecretVersion({name: secretName});
  return version.payload.data;
}

async function createTcpPool(config) {
  const pair = process.env.DB_HOST.split(':');
  return Knex({
    client: 'pg',
    connection: {
      user: process.env.DB_USER, 
      password: process.env.DB_PASS, 
      database: process.env.DB_NAME, 
      host: pair[0],
      port: pair[1]
    },
    ...config,
  });  
}

async function createUnixSocketPool(config) {
  const dbSocketPath = process.env.DB_SOCKET_PATH || '/cloudsql';
  return Knex({
    client: 'pg',
    connection: {
      user: process.env.DB_USER, 
      password: process.env.DB_PASS, 
      database: process.env.DB_NAME, 
      host: `${dbSocketPath}/${process.env.INSTANCE_CONNECTION_NAME}`,
    },
    ...config,
  });  
}

async function createPool() {
  const config = {pool: {}};
  config.pool.max = 5;
  config.pool.min = 3;
  config.pool.acquireTimeoutMillis = 60000; // 60 seconds
  config.pool.createTimeoutMillis = 30000; // 30 seconds
  config.pool.idleTimeoutMillis = 600000; // 10 minutes
  config.pool.createRetryIntervalMillis = 200; // 0.2 seconds
  const {CLOUD_SQL_CREDENTIALS_SECRET} = process.env;
  if (CLOUD_SQL_CREDENTIALS_SECRET) {
    const secrets = await accessSecretVersion(CLOUD_SQL_CREDENTIALS_SECRET);
    try {
      process.env.DB_PASS = secrets.toString();
    } catch (err) {
      err.message = `Unable to parse secret from Secret Manager. Make sure that the secret is JSON formatted: \n ${err.message} `;
      throw err;
    }
  }
  if (process.platform === 'darwin') {
    return createTcpPool(config);
  }
  return createUnixSocketPool(config);
}  

//@TODO Introduce proper migration
async function ensureSchema(pool) {
  const hasTable = await pool.schema.hasTable('stackEvent');
  if (!hasTable) {
      return pool.schema.createTable('stackEvent', table => {
      table.increments('id').primary();
      table.timestamps(true, true);
      table.string('stackId').notNullable();
      table.string('name').notNullable();
      table.string('status').notNullable();
      table.string('initiator').notNullable();
      table.string('project').notNullable();
      table.string('gcpUserAccount');
      });
  }
  console.info("Ensured that table 'stackEvent' exists");
}

let poolValue;
async function pool() {
  if (poolValue) {
    return poolValue;
  }
  poolValue = await createPool();
  await ensureSchema(poolValue);
  return poolValue;
}

let sharedSecretValue;
async function sharedSecret() {
  if (sharedSecretValue) {
    return sharedSecretValue;
  }
  const {SHARED_SECRET} = process.env;
  if (SHARED_SECRET) {
    const secrets = await accessSecretVersion(SHARED_SECRET);
    try {
      sharedSecretValue = secrets.toString();
      return sharedSecretValue;
    } catch (err) {
      err.message = `Unable to parse secret from Secret Manager. Make sure that the secret is JSON formatted: \n ${err.message} `;
      throw err;
    }
  }
}

exports.event = async (req, res) => {
  switch (req.method) {
    case 'POST':
      const key = await sharedSecret();
      if (key) {
        if (req.headers.authorization !== key) {
          res
            .status(403)
            .send("")
          return
        }
      }
      const payload = req.body
      if (!payload.stackId || !payload.name || !payload.status
         || !payload.initiator || !payload.project) {
        res
          .status(400)
          .send("")
        return
      }
      const event = {
        stackId: payload.stackId,
        name: payload.name,
        status: payload.status,
        initiator: payload.initiator,
        project: payload.project,
        gcpUserAccount: payload.gcpUserAccount,
        created_at: new Date(),
        updated_at: new Date(),
      }
      try {
        const db = await pool();
        await db('stackEvent').insert(event);
      } catch (err) {
        console.error(err)
        res
          .status(503)
          .send(err)
        return;
      }
      res
        .status(200)
        .send("")
      break;
    default:
      res
        .status(400)
        .send("Sorry, not supported")   
  }
}
