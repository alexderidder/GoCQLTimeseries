bash cqlsh -e "CREATE KEYSPACE crownstone
  WITH REPLICATION = { 
   'class' : 'SimpleStrategy', 
   'replication_factor' : 1 
  };" > ~/output
  
bash cqlsh -e "CREATE TABLE crownstone.w_and_pw_by_id_and_time ( 
    id UUID,
    time timestamp,
    w float,
	pf float
, PRIMARY KEY ((id), time)	
);" > ~/output
  
bash cqlsh -e "CREATE TABLE crownstone.kwh_by_id_and_time ( 
    id UUID,
    time timestamp,
    kwh float,
, PRIMARY KEY ((id), time)	
);" > ~/output