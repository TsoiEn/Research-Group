document.addEventListener('DOMContentLoaded', () => {
    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');

    function showTab(tabId) {
        tabButtons.forEach(button => button.classList.remove('active'));
        tabContents.forEach(content => content.classList.remove('active'));

        document.querySelector(`[onclick="showTab('${tabId}')"]`).classList.add('active');
        document.getElementById(tabId).classList.add('active');
    }

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('onclick').match(/'([^']+)'/)[1];
            showTab(tabId);
        });
    });
});

function showPDFModal(base64Data) {
    const pdfViewer = document.getElementById('pdfViewer');
    pdfViewer.src = `data:application/pdf;base64,${base64Data}`;
    document.getElementById('pdfModal').style.display = 'block';
}

function closePDFModal() {
    document.getElementById('pdfModal').style.display = 'none';
}

function openAddCredentialModal() {
    document.getElementById("addCredentialModal").style.display = "block";
}

function closeAddCredentialModal() {
    document.getElementById("addCredentialModal").style.display = "none";
}
