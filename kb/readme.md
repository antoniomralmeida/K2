## O Projeto K2  é um projeto de pesquisa científica e inovação que objetiva a construção de um sistema especialista que utiliza recursos tecnológicos de inteligência artificial, bem como outras técnicas avançada de computação. O objetivo deste sistema é apoiar o trabalho de um especialista na tomada de decisões táticas e operacionais, quando o recurso especialista for caro e de difícil obtenção e retenção ou quando o trabalho envolve ter que tomar várias decisões complexas em um tempo muito curto ou ainda quando o ambiente de trabalho do especialista for insalubre.


### Os produtos do projeto K2 podem ser usados em várias áreas de trabalho, como:
1.	Gerenciar plantas industriais de forma autônoma (PLCs)
2.	Reduzir tempo de setup de linhas industriais
3.	Gerenciar centrais de atendimento ou de gerenciamento de serviços
4.	Definir de forma autônoma limite de crédito de clientes
5.	Definir de forma autônoma preço de produtos e serviços
6.	Telemedicina, com diagnóstico médico complementar
7.	Uso militar no controle de lançamento de foguetes e misseis, teleguiados
8.	Na gestão de risco de qualquer natureza, inclusive ambiental


### TODO List

└─ K2
   ├─ kb
   │  ├─ kbattributeobject.go
   │  │  ├─ line 16: TODO : Acionar regras em backward
   │  │  ├─ line 17: TODO : Create uma tarefa de simulação
   │  │  ├─ line 18: TODO : As tarefas de busca de valor devem ter timite de tempo
   │  │  ├─ line 19: TODO : Criar formulário web para receber valores de atributos de origem User (assincrono)
   │  │  ├─ line 20: TODO : Levar em consideração a certeze na obteção de um valor PLC e User 100%, criar regra de envelhecimento da certeza
   │  │  ├─ line 21: TODO : a certeza de um valor simulado deve analizer os quadrantes da curva normal do historico de valor
   │  │  └─ line 22: TODO : a certeza por inferencia deve usar logica fuzzi
   │  ├─ kbclass.go
   │  │  └─ line 39: TODO : Restart KB
   │  ├─ kbobject.go
   │  │  └─ line 22: TODO : Reiniciar KB
   │  ├─ kbrule.go
   │  │  ├─ line 77: FIXME : A classe está vazia e não deveria
   │  │  └─ line 205: TODO : Acionar regras em forward chaining
   │  └─ readme.md
   │     └─ line 15: TODO List
   └─ web
      ├─ bootstrap
      │  ├─ js
      │  │  ├─ bootstrap.bundle.js
      │  │  │  └─ line 6382: TODO (fat): remove sketch reliance on jQuery position/offset
      │  │  ├─ bootstrap.bundle.js.map
      │  │  │  └─ line 1:  
      │  │  ├─ bootstrap.bundle.min.js.map
      │  │  │  └─ line 1:  
      │  │  ├─ bootstrap.js
      │  │  │  └─ line 3769: TODO (fat): remove sketch reliance on jQuery position/offset
      │  │  ├─ bootstrap.js.map
      │  │  │  └─ line 1:  
      │  │  └─ bootstrap.min.js.map
      │  │     └─ line 1:  
      │  └─ scss
      │     └─ _reboot.scss
      │        └─ line 33: TODO : remove in v5
      ├─ chart.js
      │  ├─ Chart.bundle.js
      │  │  ├─ line 2996: TODO (v3): remove 'global' from namespace.  all default are global and
      │  │  ├─ line 9547: TODO (SB): I think we should be able to remove this custom case (options.scale)
      │  │  ├─ line 11593: TODO (v3): remove minSize as a public property and return value from all layout boxes. It is unused
      │  │  ├─ line 13148: TODO (v3): change this to positiveOrDefault
      │  │  ├─ line 15067: TODO : Remove "ordinalParse" fallback in next major release.
      │  │  ├─ line 15710: TODO : add sorting
      │  │  ├─ line 15750: TODO : Another silent failure?
      │  │  ├─ line 16555: TODO : Find a better way to register and load all the locales in Node
      │  │  ├─ line 16835: TODO : We need to take the current isoWeekYear, but that depends on
      │  │  ├─ line 17010: TODO : Replace the vanilla JS Date object with an indepentent day-of-week check.
      │  │  ├─ line 17114: TODO : Move this to another part of the creation flow to prevent circular deps
      │  │  ├─ line 17409: TODO : Use [].sort instead?
      │  │  ├─ line 17840: TODO : remove 'name' arg after deprecation is removed
      │  │  ├─ line 18468: TODO : Remove "ordinalParse" fallback in next major release.
      │  │  └─ line 19000: TODO : Use this.as('ms')?
      │  └─ Chart.js
      │     ├─ line 2992: TODO (v3): remove 'global' from namespace.  all default are global and
      │     ├─ line 9543: TODO (SB): I think we should be able to remove this custom case (options.scale)
      │     ├─ line 11589: TODO (v3): remove minSize as a public property and return value from all layout boxes. It is unused
      │     └─ line 13144: TODO (v3): change this to positiveOrDefault
      ├─ fontawesome-free
      │  ├─ all.js
      │  │  └─ line 2186: TODO : do we need to handle font-weight for kit SVG pseudo-elements?
      │  └─ fontawesome.js
      │     └─ line 202: TODO : do we need to handle font-weight for kit SVG pseudo-elements?
      └─ jquery
         ├─ jquery.js
         │  ├─ line 794: TODO : identify versions
         │  ├─ line 808: TODO : identify versions
         │  └─ line 4466: TODO : Now that all calls to _data and _removeData have been replaced
         └─ jquery.slim.js
            ├─ line 794: TODO : identify versions
            ├─ line 808: TODO : identify versions
            └─ line 4466: TODO : Now that all calls to _data and _removeData have been replaced
