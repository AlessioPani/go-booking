# go-booking
This repository contains a web application for a Booking and Reservation system.

- Built in Go version 1.19
- Uses the following external dependences 
  - [chi router](https://github.com/go-chi/chi/v5)
  - [Alex Edwards's SCS session management system](https://github.com/alexedwards/scs/v2)
  - [NoSurf](https://github.com/justinas/nosurf) as a middleware
  - [Soda/Pop](https://gobuffalo.io/documentation/database/pop/) for database migrations
  - [Alex Saskevich's](https://github.com/asaskevich/govalidator) validator
  - [Pgx](https://github.com/jackc/pgx) as driver for the PostgreSQL database
  - [Godotenv](https://github.com/joho/godotenv) to load a .env file with potentially must-remain-secret data
  
- To run

  - ```shell
    ./run.sh
    ```


* Testing

  * ```shell
    go test (-v)
    ```

  * ```shell
    go test -coverprofile=coverage.out && go tool cover -html=coverage.out
    ```

  
  - ```shell
    go test -v ./...
    ```
  

- Various used instructions

  - ```shell
    soda generate fizz <name_of_migration>
    ```

  - ```
    soda migrate
    ```

  - ```
    soda migrate down
    ```

    

## Database structure

### User

Table used for login and authentication of the website owner to check on the backend, with the following fields:

- id
- first_name
- last_name
- email
- password
- created_at (automatically created by soda)
- updated_at (automatically created by soda)
- access_level

### Rooms

Table used to save the information of each room, with the following fields:

- id
- room_name
- created_at (automatically created by soda)
- updated_at (automatically created by soda)

### Reservations

Table used to hold all the details of a single reservation, with the following fields:

- id
- first_name
- last_name
- email
- phone
- start_date
- end_date
- room_id (foreign key to table Rooms)
- created_at (automatically created by soda)
- updated_at (automatically created by soda)

### Restrictions

Table used to save a list of restriction options, with the following fields:

- id
- restriction_name
- created_at (automatically created by soda)
- updated_at (automatically created by soda)

### Room Restrictions

Table used to hold the information about restrictions over a room in a given period of time with a specific reason determined by restriction_id, with the following fields: 

- id
- start_date
- end_date
- room_id (foreign key to table Rooms)
- restriction_id (foreign key to table Restrictions)
- reservation_id (foreign key to table Reservations)
- created_at (automatically created by soda)
- updated_at (automatically created by soda)

