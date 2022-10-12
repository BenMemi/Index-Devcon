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


There is also a dockerfile if you want to deploy your indexer/run it in a compose 

You can just make production PROJECT_ID=YOUR_GCP_PROJECT_ID and it will push up your container to the registry for you! 

Also you can make test to test the container locally before pushing it up to the registry!
      