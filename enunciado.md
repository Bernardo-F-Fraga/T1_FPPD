Ir para o conteúdo principal
Condições de conclusão

Implementar o algoritmo distribuído de eleição baseado em anel lógico (Ring Algorithm) utilizando o código base fornecido na linguagem Go:

    Baseado no uso de um anel lógico;
    Cada processo conhece o anel inteiro, mas manda mensagens somente para o próximo processo ativo no sentido do anel;
    Quando processo detecta que o coordenador não está mais ativo (isto pode ser simulado com uma mensagem de um processo externo), ele envia uma mensagem de eleição no anel contendo seu id (disputando esta coordenação que está vaga);
    A cada passo, o processo que recebe a mensagem de eleição inclui seu id na mensagem e passa adiante no anel
    No final de uma volta completa, o processo que iniciou a eleição recebe a mensagem e escolhe aquele que tem maior id como novo processo coordenador;
    Uma nova mensagem é enviada através do anel para que todos conheçam o novo coordenador;

Modelar uma forma de mostrar que está funcionando (ex: processo externo que simula solicitações) e também de simular falhas para que a resposta não seja sempre a mesma (nodos do anel que param de funcionar).

A avaliação do trabalho será feita com base no acompanhamento do desenvolvimento em laboratório e da entrega de um relatório.

Formato do relatório técnico:

- arquivo formato .pdf;
- cabeçalho reduzido com identificação do grupo e do trabalho;
- primeira página coluna dupla com margens reduzidas (1cm);

- segunda página com dumps de tela mostrando interface e mensagens do programa em funcionamento (referenciar figuras no texto);
- a partir da terceira página código fonte formatado em coluna simples (sem limite). Numerar as linhas de código para possibilitar referências no texto

Relatório deve incluir:

- breve introdução do problema da eleição e pra que é usado em sistemas distribuídos citando um exemplo prático

- breve descrição do algoritmo específico de eleição que está sendo implementado (ring) e suas fases

- descrição dos tipos de mensagens foram usadas na implementação e como foram usadas

- descrição de como foram modeladas as falhas

- apresentação de como o algoritmo foi testado e que operações foram feitas nestes testes (papel do processo externo)

(ver a rubrica de avaliação na sala de envio para mais detalhes)

Última atualização: quarta-feira, 29 abr. 2026, 15:37
