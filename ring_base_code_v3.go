// Código exemplo para o trabaho de sistemas distribuidos (eleicao em anel)
// By Cesar De Rose - 2026

package main

import (
	"fmt"
	"sync"
)

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [4]int // conteudo da mensagem para colocar os ids (usar um tamanho compativel com o numero de processos no anel)
}

var (
	chans = []chan mensagem{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
	}
	controle = make(chan int)
	wg       sync.WaitGroup // wg is used to wait for the program to finish

	// IMPORTANTE - em sistemas distribuídos não temos memória compartilha de forma que 
	// não podem ser usadas variáveis globais no trbalho além das que foram declaradas acima
)

func ElectionControler(in chan int) {
	defer wg.Done()

	var temp mensagem

	// comandos para o anel iniciam aqui (para servir de exemplo)

	// MK: Mata processo 3 (entrada pelo canal 2)
	temp.tipo = 2
	chans[2] <- temp
	fmt.Printf("	Controle: mudar o processo 3 para falho\n")

	fmt.Printf("	Controle: confirmação recebida do processo %d\n", <-in) // receber e imprimir confirmação

	// MK: Solicita eleição ao processo 2 (entrada pelo canal 1)
	temp.tipo = 1
	fmt.Printf("	Controle: solicitar eleição ao processo 2\n")
	chans[1] <- temp

	fmt.Printf("	Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

		fmt.Printf("\n	Controle: ELEIÇÃO FINALIZADA \n\n")

	// MK: Mata processo 2 (entrada pelo canal 1)
	temp.tipo = 2
	chans[1] <- temp
	fmt.Printf("	Controle: mudar o processo 2 para falho\n")

	fmt.Printf("	Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// MK: Solicita eleição ao processo 1 (entrada pelo canal 0)
	temp.tipo = 1
	fmt.Printf("	Controle: solicitar eleição ao processo 1\n")
	chans[0] <- temp

	fmt.Printf("	Controle: confirmação recebida do processo %d\n", <-in) // receber e imprimir confirmação

		fmt.Printf("\n	Controle: ELEIÇÃO FINALIZADA \n\n")

	// MK: Mata processo 1 (entrada pelo canal 0)
	temp.tipo = 2
	chans[0] <- temp
	fmt.Printf("	Controle: mudar o processo 1 para falho\n")

	fmt.Printf("	Controle: confirmação recebida do processo %d\n", <-in) // receber e imprimir confirmação
	
	// MK: Solicita Eleição ao processo 0 (entrada pelo canal 3)
	temp.tipo = 1
	fmt.Printf("	Controle: solicitar eleição ao processo 0\n")
	chans[3] <- temp

	fmt.Printf("	Controle: confirmação recebida do processo %d\n", <-in) // receber e imprimir confirmação

	fmt.Printf("\n	Controle: ELEIÇÃO FINALIZADA \n\n")

	// MK: Revive processo 3 (entrada pelo canal 2)
	temp.tipo = 3
	chans[2] <- temp
	fmt.Printf("	Controle: mudar o processo 3 para vivo\n")

	fmt.Printf("	Controle: confirmação recebida do processo %d\n", <-in) // receber e imprimir confirmação

	// MK: Solicita Eleição ao processo 3 (entrada pelo canal 2)
	temp.tipo = 1
	fmt.Printf("	Controle: solicitar eleição ao processo 3\n")
	chans[2] <- temp

	fmt.Printf("	Controle: confirmação recebida do processo %d\n", <-in) // receber e imprimir confirmação

	fmt.Printf("\n	Controle: ELEIÇÃO FINALIZADA \n\n")


	temp.tipo = 4
	chans[0] <- temp
	chans[1] <- temp
	chans[2] <- temp
	chans[3] <- temp

	fmt.Println("\n   Processo controlador concluído\n")
}

func TaskProcess(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	// variaveis locais que indicam se este processo é o lider e se esta ativo
	// estas variáveis só podem ser alteradas pelo recebimento de mensagens nos canais

	var actualLeader int
	var bFailed bool = false // todos inciam sem falha
	actualLeader = leader // indicação do lider veio por parâmatro

	
	var finish bool = false

	for !finish {
		temp := <-in
					
		if !bFailed {
			fmt.Printf("%2d: recebi mensagem: %d, [ %d, %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])
		} else {
			fmt.Printf("%2d: (FALHO) recebi mensagem: %d, [ %d, %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])
		}
		switch temp.tipo {
		case 0:// MK: Recebeu mensagem de eleição do processo anterior
			{
				if bFailed {
					out <- temp
				} else {
					temp.corpo[TaskId] = TaskId
					out <- temp
				}
			}
		case 1: // MK: Controlador Solicitou Eleição
			{
				fmt.Printf("%2d: Convocando eleição\n", TaskId)
				temp.tipo = 0
				for i := range temp.corpo {
        			temp.corpo[i] = -5 // MK: Utilizei -5 como valor nulo, não definido
    			}
				temp.corpo[TaskId] = TaskId
				out <- temp // MK: Envia mensagem para o próximo no anel solicitando votação.

				temp = <-in // MK: Aguarda travado a mensagem dar a volta completa no anel e retornar para cá novamente.
				fmt.Printf("%2d: votação deu a volta no anel %d, [ %d, %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])

				actualLeader = -5
				for i := range temp.corpo {
					if temp.corpo[i] > actualLeader {
						actualLeader = temp.corpo[i]
					}
				}

				temp.corpo[0] = actualLeader
				temp.tipo = 5
				out <- temp

				temp = <- in
				
				fmt.Printf("%2d: mudou líder atual para %d\n", TaskId, actualLeader)
				controle <- TaskId
			}
		case 2: // MK: Processo Falhou
			{
				bFailed = true
				fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
				controle <- TaskId
			}
		case 3: // MK: Processo Reviveu
			{
				bFailed = false
				fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
				controle <- TaskId
			}
		case 4: // MK: Processo Terminou
			{
				finish = true
			}
		case 5: // MK: Eleição concluída - Recebendo informação de novo líder
			{
				if !bFailed {
					actualLeader = temp.corpo[0]
					fmt.Printf("%2d: Mudou líder atual %d\n", TaskId, actualLeader)
				}
				out <- temp
			}
		default: // MK: Mensagem não Reconhecida
			{
				fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
			}
		}
	}
	fmt.Printf("%2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // Add a count of four, one for each goroutine

	// criar os processo do anel de eleicao

	go TaskProcess(0, chans[3], chans[0], 3) // este é o lider
	go TaskProcess(1, chans[0], chans[1], 3) // não é lider, é o processo 0
	go TaskProcess(2, chans[1], chans[2], 3) // não é lider, é o processo 0
	go TaskProcess(3, chans[2], chans[3], 3) // não é lider, é o processo 0

	fmt.Println("\n   Anel de processos criado")

	// criar o processo controlador

	go ElectionControler(controle)

	fmt.Println("\n   Processo controlador criado\n")

	wg.Wait() // Wait for the goroutines to finish\
}