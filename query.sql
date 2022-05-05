select * from "stackEvent" order by created_at DESC;

select count(*) from "stackEvent" where status = 'deployed';
select count(*) from "stackEvent" where status = 'undeployed';
select count(*) from "stackEvent" where status = 'incomplete' 
    and project = 'superhub';

select project, count(*) from "stackEvent" 
    where status = 'deployed' 
    group by project;

select initiator, project, count(*) from "stackEvent"
    where status = 'deployed'
    group by initiator, project

select project, name, status, count(*) from "stackEvent" 
    where status = 'deployed' or status = 'undeployed' or status = 'incomplete'
    group by project, name, status;

select name, count(*) from "stackEvent" 
    where status = 'incomplete' 
    group by name;
    