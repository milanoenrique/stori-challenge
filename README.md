# stori-challenge
In this project you can find a a code for run a job and process csv file into the transctions folder,the of every file must be the account id of the client, is necessary configure the db password into docker-compose.yml and the gmail account for send the email with the account resume to the users.

The job execute 5 workers for process in paralell all files into transactions folders and send an email with the account information to the client.

For run the project compile the project with docker-compose build --up, and later make a cron job with the sintax: docker run --rm my-job-image:latest
