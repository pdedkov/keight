# keight

## How to use

### Run app
1. cp env.example .env.docker (configure app port)
2. cp docker-compose.yml.exampl docker-compose.yml
3. docker-compose up -d api

### Test upload
 - upload to tests dir file or use file.tar.gz (docker volume)
 - replace TEST_UPL_FILE in docker-compose.yml (file in)
 - docker-compose up upl (replace in docker-composer.yml)

### Test download
 - replace UPLOAD_ID with id recieved in result of upload request
 - docker-compose up dl (file will download to tests dir (docker volume))

## Limitations
1. Upload 1 file per request
2. After restart storage all data in storage will erase
3. every upload locks processing because of dummy storage
4. every download locks processing because of dummy storage