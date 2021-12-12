# CurrencyChecker

You can run the project using

make run

After running the project, you can send get requests to

localhost:9999/{Currency}

Example request

curl --location --request GET 'localhost:9999/EUR'

Example response

{"data":8.5,"currency":"EUR"}