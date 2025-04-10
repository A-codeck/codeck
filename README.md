# CODECK
###### COmpetição DE Código do Kim 
GymRats para maratoners. Codar é legal, mas codar melhor que seus amigos é mais.

# Banco de dados
Utilizamos PostgreSQL para o banco de dados. 
## Diagrama de tabelas
```mermaid
%%{init: {'theme': 'forest', 'themeVariables': {
    'primaryColor': '#ffcc00',
    'edgeLabelBackground':'#ffffff',
    'tertiaryColor': '#e0e0e0',
    'lineColor':'#ffffff'
}}}%%

erDiagram
    %% Tables Definition
    USERS {
      integer id PK
      varchar email
      varchar nome
      varchar infos_de_login
    }
    ACTIVITIES {
      integer id PK
      integer user_id
      varchar titulo
      varchar desc
      date data
      varchar imagem
    }
    COMMENTS {
      integer id PK
      integer activity_id
      integer group_id
      varchar text
      date data
    }
    GROUPS {
      integer id PK
      varchar nome
      varchar desc
      date data_inicio
      date data_fim
    }
    USER_GROUPS {
      varchar apelido
      integer group_id PK
      integer user_id PK
    }
    GROUP_ACTIVITIES {
      integer acitivity_id PK
      integer group_id PK
    }

    %% Relationships (Foreign Keys)
    USERS ||--o{ ACTIVITIES : user_id
    USERS ||--o{ USER_GROUPS : user_id

    ACTIVITIES ||--o{ COMMENTS : activity_id
    ACTIVITIES ||--o{ GROUP_ACTIVITIES : acitivity_id

    GROUPS ||--o{ COMMENTS : group_id
    GROUPS ||--o{ GROUP_ACTIVITIES : group_id
    GROUPS ||--o{ USER_GROUPS : group_id

```
# Diagrama de classes 
```mermaid
%%{init: {'theme': 'forest', 'themeVariables': {
    'primaryColor': '#ffcc00',
    'edgeLabelBackground':'#ffffff',
    'tertiaryColor': '#e0e0e0',
    'lineColor':'#ffffff'
}}}%%

classDiagram
    %% Models
    class UserModel {
      +int id
      +String name
      +String email
      +List~GroupModel~ groups
    }
    class GroupModel {
      +int id
      +String name
      +List~UserModel~ members
      +List~ActivityModel~ activities
    }
    class ActivityModel {
      +int id
      +String title
      +String description
      +Date date
      +List~CommentModel~ comments
    }
    class CommentModel {
      +int id
      +String text
      +Date timestamp
    }

    %% Relationships between models
    UserModel "1" --> "0..*" GroupModel 
    GroupModel "1" --> "0..*" ActivityModel 
    ActivityModel "1" --> "0..*" CommentModel 

    %% Controllers
    class UserController {
      +getUserInfo()
      +accessProfile()
      +listActivities()
      +manageGroupRoles()
    }
    class LoginController {
      +createAccount()
      +loginUser()
    }
    class GroupController {
      +createGroup()
      +readGroup()
      +editGroup()
      +deleteGroup()
      +addUserToGroup()
      +listGroupActivities()
    }
    class ActivityController {
      +createActivity()
      +readActivity()
      +updateActivity()
      +listActivityComments()
    }
    class CommentController {
      +createComment()
    }

    %% Controller uses model
    UserController --> UserModel
    LoginController --> UserModel
    GroupController --> GroupModel
    GroupController --> UserModel
    ActivityController --> ActivityModel
    ActivityController --> CommentModel
    CommentController --> CommentModel

    %% Views
    class Homepage {
      +loginButton
    }
    class LoginView {
      +enterCredentials()
      +createAccountLink()
    }
    class HomeView {
      +displayRecentActivities()
      +displayCalendar()
      +displayGroupList()
    }
    class GroupView {
      +displayGroupActivities()
      +displayProgressBar()
      +displayGroupRanking()
    }
    class ActivityCreationView {
      +displayActivityForm()
    }
    class ActivityDetailView {
      +displayActivityDetails()
    }

    %% Views interact with controllers
    Homepage --> LoginController
    LoginView --> LoginController
    HomeView --> ActivityController
    HomeView --> GroupController
    GroupView --> GroupController
    ActivityCreationView --> ActivityController
    ActivityDetailView --> ActivityController

```
