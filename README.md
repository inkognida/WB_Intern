# WB_Intern

Run: 
  - docker-compose up -d
  - cd cmd/
  - go run main.go

To test it: 
  - cd cmd/testpusher/
  - go run main.go

Check: 
  - localhost:8080
  - docker exec -it {container_id} psql -U admin -d orders
