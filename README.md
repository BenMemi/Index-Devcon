# SIMPLE INDEXER 

Just look through the comments in main.go and adjust where specified to have your own indexer working for your events :) 

Also define your db tables in the database directory/module! 

the .env should look something like this; 

DATABASE_URL=user=postgres password=postgres host=db.postgres.supabase.co port=5432 dbname=postgres
RPC_URL=wss://polygon-mumbai.g.alchemy.com/v2/-no-touchy
HISTORIC = false
FROM_BLOCK=1
TO_BLOCK=5
RPC_HISTORIC=https://polygon-mumbai.gateway.pokt.network/v1/lb/no-touchy

      