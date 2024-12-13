
# Postgres Password Generator

## Description

Postgres Password Generator is a Go application that allows you to generate complex passwords, securely store them in a PostgreSQL database, and manage them using console flags. The application ensures password encryption during storage and provides secure interaction with the database.

## Features

- **Random Password Generation**: The application generates complex and random passwords that are hard to guess.
- **Password Encryption**: Passwords are encrypted before storage in the database for security.
- **Password Management**: The application allows you to add, delete, and search for passwords using convenient console flags.
- **Docker Container**: The application launches and manages a Docker container with a PostgreSQL database, making installation and configuration easy.

## Technologies Used
- **Go**: The programming language used for the application development.
- **PostgreSQL**: The database management system used for storing passwords.
- **Docker**: Used to launch and manage a container with a PostgreSQL database.
- **AES-GCM Encryption**: Used for password encryption and decryption.

## Installation and Launch

1. **Install Docker**:
   - Ensure that Docker is installed on your system. If not, install it from the official website: [Docker](https://www.docker.com/).

2. **Clone the Repository**:
   - Clone this repository to your system:
     ```
     git clone https://github.com/your-username/password-generator.git
     cd password-generator
     ```

3. **Set Environment Variables**:
   - Create an `.env` file in the project root and add the following variables:
     ```
     CONTAINER_NAME=your_container_name
     CONTAINER_USER=postgres_user
     DB_PASSWORD=your_db_password
     DB_PORT=5432
     DB_NAME=postgres
     DB_HOST=localhost
     SSL_MODE=disable
     HOST_PORT=4000
     ENCRYPTION_KEY=lgjsgwioghnbtyib
     ```
     Replace `your_container_name`, `your_db_password`, `ENCRYPTION_KEY` value and other values with the appropriate ones for your system.

     **Important Note**: The `ENCRYPTION_KEY` in the `.env` file should be a string of exactly 16 ASCII characters to ensure 16-byte encryption key length.

4. **Build the Application**:
   - On Unix-based systems, you can use `make` to build the application:
     ```
     make build
     ```
     This will create a binary named `pass_gen` in the project directory.

   - **For Windows Users**: If you are using Windows, run the following command instead:
     ```
     go build -o pass_gen
     ```

5. **Run the Application**:
   - On Unix-based systems:
     ```
     make dev
     ```
     This will build and launch the application in one step.

   - **For Windows Users**: Run the application manually:
     ```
     ./pass_gen
     ```

## Usage

The application provides various console flags for password management:

- **`-showme`**: Displays the content of the passwords table.
- **`-drop`**: Clears the passwords table.
- **`-delete`**: Deletes a password by its ID.
- **`-delete-last`**: Deletes multiple recent passwords from the table.
- **`-find-service`**: Displays passwords associated with a specified service.
- **`-find-password`**: Displays a password by its ID.
- **`-no-save`**: Outputs a new password without saving it in the database.

If no flags are specified, the application generates a new random password, encrypts it, and stores it in the database.

## Project Structure

The project consists of several packages, each responsible for specific functionality:

- **`main`**: The main package containing the application's entry point.
- **`container_manager`**: Manages the Docker container with the PostgreSQL database.
- **`flag_manager`**: Implements logic for handling console flags for password management.
- **`repository`**: Contains functions for interacting with the database.
- **`cypher`**: Includes functions that implement encryption and decryption of passwords.
- **`package/environment`**: Provides functions for working with environment variables.
- **`package/password_generator`**: Generates random passwords.

## Security

Security is a priority in this project. The application ensures the following security measures:

- **Password Encryption**: Passwords are encrypted using the AES-GCM algorithm, ensuring their confidentiality and integrity.
- **Secure Storage of Encryption Key**: The encryption key should be stored in a secure location, such as the `.env` file, and should not be included in the repository.
- **Password Existence Check**: The application checks for the existence of a password in the database before adding a new one to prevent duplicates.

## Other

The project contains commented-out //NOTE console logs in order to facilitate possible application debugging

___

# Postgres Password Generator

## Описание
Postgres Password Generator - приложение, написанное на Go, которое позволяет генерировать сложные пароли, безопасно хранить их в базе данных PostgreSQL, а также управлять ими с помощью консольных флагов. Приложение обеспечивает шифрование паролей при хранении и обеспечивает безопасное взаимодействие с базой данных.

## Особенности
- **Генерация случайных паролей**: приложение генерирует сложные и случайные пароли, которые сложно угадать.
- **Шифрование паролей**: пароли шифруются перед хранением в базе данных для обеспечения безопасности.
- **Управление паролями**: приложение позволяет добавлять, удалять и искать пароли с помощью удобных консольных флагов.
- **Docker-контейнер**: приложение запускает и управляет контейнером Docker с базой данных PostgreSQL, обеспечивая простоту установки и настройки.

## Используемые технологии
- **Go**: Язык программирования, используемый для разработки приложений.
- **PostgreSQL**: Система управления базами данных, используемая для хранения паролей.
- **Docker**: Используется для запуска контейнера с базой данных PostgreSQL и управления им.
- **Шифрование AES-GCM**: Современный алгоритм шифрования, обеспечивающий безопасность паролей.

## Установка и запуск

1. **Установите Docker**:
   - Убедитесь, что Docker установлен на вашей системе. Если его нет, установите его с официального сайта: [Docker](https://www.docker.com/).

2. **Клонируйте репозиторий**:
   - Склонируйте этот репозиторий на вашу систему:
     ```
     git clone https://github.com/your-username/password-generator.git
     cd password-generator
     ```

3. **Настройте переменные окружения**:
   - Создайте файл `.env` в корне проекта и добавьте в него следующие переменные:
     ```
     CONTAINER_NAME=your_container_name
     CONTAINER_USER=postgres_user
     DB_PASSWORD=your_db_password
     DB_PORT=5432
     DB_NAME=postgres
     DB_HOST=localhost
     SSL_MODE=disable
     HOST_PORT=4000
     ENCRYPTION_KEY=lgjsgwioghnbtyib
     ```
     Замените `your_container_name`, `your_db_password`, значение `ENCRYPTION_KEY` и другие значения на подходящие для вашей системы.

     **Важно**: Переменная `ENCRYPTION_KEY` в файле `.env` должна быть строкой **ровно из 16 символов ASCII**, чтобы обеспечить длину ключа шифрования в 16 байт.

4. **Соберите приложение**:
   - На Unix-подобных системах вы можете использовать `make` для сборки приложения:
     ```
     make build
     ```
     Эта команда создаст бинарный файл с именем `pass_gen` в корневой директории проекта.

   - **Для пользователей Windows**: Если вы используете Windows, выполните следующую команду вместо `make`:
     ```
     go build -o pass_gen
     ```

5. **Запустите приложение**:
   - На Unix-подобных системах:
     ```
     make dev
     ```
     Эта команда одновременно соберет и запустит приложение.

   - **Для пользователей Windows**: Запустите приложение вручную:
     ```
     ./pass_gen
     ```

## Использование

Приложение предоставляет различные консольные флаги для управления паролями:

`-showme`: выводит содержимое таблицы паролей.
`-drop`: очищает таблицу паролей.
`-delete`: удаляет пароль по его ID.
`-delete-last`: удаляет несколько последних паролей из таблицы.
`-find-service`: выводит пароли, связанные с указанной службой.
`-find-password`: выводит пароль по его ID.
- **`-no-save`**: Выводит новый пароль без сохранения его в базе данных.

Если ни один из флагов не указан, приложение генерирует новый случайный пароль, шифрует его и сохраняет в базе данных.

## Структура проекта
Проект с��стоит из нескольких пакетов, каждый из которых отвечает за определенную функциональность:

**`main`**: основной пакет, содержащий точку входа в приложение.
**`container_manager`**: отвечает за управление контейнером Docker с базой данных PostgreSQL.
**`flag_manager`**: реализует логику обработки консольных флагов для управления паролями.
**`repository`**: содержит функции для взаимодействия с базой данных.
**`cypher`**: включает в себя функции, реализующие шифрование и расшифровку паролей.
**`package/environment`**: предоставляет функции для работы с переменными окружения.
**`package/password_generator`**: генерирует случайные пароли.

## Безопасность
Безопасность является приоритетом в этом проекте. Приложение обеспечивает следующие меры безопасности:

**Шифрование паролей**: пароли шифруются с помощью алгоритма AES-GCM, что гарантирует их конфиденциальность и целостность.
**Безопасное хранение ключа шифрования**: ключ шифрования должен храниться в безопасном месте, таком как файл `.env`, и не включаться в репозиторий.
**Проверка существования пароля**: приложение проверяет существование пароля в базе данных перед добавлением нового, чтобы избежать дублирования.

## Другое 

В проекте присутствуют закомментированные //NOTE выводы в консоль для облегчения возможного дебага приложения 
