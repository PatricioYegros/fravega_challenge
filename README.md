# fravega_challenge
API REST called Ordes responsible for managing customer purchase orders

- How to execute the API locally:

    You will need:
        Docker
        Postman
        CMD or another command terminal.

    1) With your command terminal go to the directory: challenge_pyegros/app. (Ensure you have docker running!)

    2) There you can execute the follow command "docker compose up -d".
        This command will create the images of the API and the DataBase in a Docker Container.
    
    3) Import the "Fravega.postman_collection.json" to your Postman

    3) Now you are ready to go!

- Technical Decisions:

    1) When you use the Docker container for the first time, you will need to create a order first.
        The database is empty at that moment.

    2) The autoincremental id is managed manually by the code but automatic for the user.
    
    3) The translations are done by two simply switch because its easier and faster than implementing i18n or CGP translate.
        If the translations grow in quantity or languages, i would implement i18n or if i have a project in GCP, the service
        of Google Translate.

    4) The correct date format is checked (is RFC3339).

    5) The idempotency is handled through Redis Caché, with a TTL of 1 day.

    6) The externalReferenceId of each Channel are made up by me and they are these:

            "Ecommerce":  "abc-123",
            "CallCenter": "def-456",
            "Store":      "ghi-789",
            "Affiliate":  "jkl-012",