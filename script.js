// Inicia no modo básico
document.addEventListener("DOMContentLoaded", function () {
    mostrarProximaPergunta();
});

function selecionarVersao(versao) {
    // Lógica para selecionar a versão (se necessário)
    console.log("Versão selecionada: " + versao);
}

function mostrarProximaPergunta() {
    // Esconde a pergunta atual
    var perguntaAtual = document.querySelector('.pergunta.visivel');
    if (perguntaAtual) {
        perguntaAtual.classList.remove('visivel');
    }

    // Mostra a próxima pergunta
    var proximaPergunta = perguntaAtual.nextElementSibling;
    if (proximaPergunta) {
        proximaPergunta.classList.add('visivel');
    } else {
        // Se não houver mais perguntas, mostra o formulário do paciente
        document.getElementById('perguntasPaciente').style.display = 'block';
    }
}

function cadastrarPaciente() {
    var nome = document.getElementById("nome").value;
    var idade = document.getElementById("idade").value;
    var email = document.getElementById("email").value;

    // Adiciona o paciente à tabela
    adicionarPacienteNaTabela(nome, idade, email);

    // Limpa os campos do formulário
    document.getElementById("nome").value = "";
    document.getElementById("idade").value = "";
    document.getElementById("email").value = "";

    // Reinicia o formulário
    reiniciarFormulario();
}

function adicionarPacienteNaTabela(nome, idade, email) {
    var tabela = document.getElementById('tabelaBody');
    var novaLinha = tabela.insertRow(tabela.rows.length);

    var colunaNome = novaLinha.insertCell(0);
    colunaNome.innerHTML = nome;

    var colunaIdade = novaLinha.insertCell(1);
    colunaIdade.innerHTML = idade;

    var colunaEmail = novaLinha.insertCell(2);
    colunaEmail.innerHTML = email;
}

function reiniciarFormulario() {
    // Esconde as perguntas do paciente
    document.getElementById('perguntasPaciente').style.display = 'none';

    // Reseta a tabela de pacientes
    var tabela = document.getElementById('tabelaBody');
    while (tabela.firstChild) {
        tabela.removeChild(tabela.firstChild);
    }

    // Volta para a primeira pergunta
    document.getElementById('pergunta1').classList.add('visivel');
}
