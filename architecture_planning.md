#Architecture planning

## Models

UserModel

GroupModel

ActivityModel

CommentModel

## Controllers
UserController
 - Ler informações do usuário
 - Acessar perfil
 - Listar atividades (?)
 - GroupRoles

LoginController
 - Criar conta
 - Logar usuário

GroupController
 - Criar grupo
 - Ler grupo
 - Editar grupo
 - Deletar grupo
 - Adicionar pessoa a grupo
 - Listar atividades do grupo (?)

ActivityController
 - Criar atividade
 - Ler atividade (detalhes da atividade)
 - Atualizar atividade
 - Ler comentários de uma atividade

CommentController
 - Criar comentário


## View
Homepage 
 - Botão para login

Tela de login 
 - Entrar na conta
 - Criar conta

Tela inicial 
 - Atividade recente em todos os grupos
 - Calendário com a sua streak de problemas
 - Lista de grupos

Tela de um grupo
 - Atividades postadas naquele grupo
 - Barra de progresso
 - Ranking do grupo

Tela criação atividade

Tela de detalhes de uma atividade

