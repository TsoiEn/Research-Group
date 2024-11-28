document.addEventListener("DOMContentLoaded", function () {
    console.log("Script Loaded");

    // Function to switch tabs
    window.showTab = function (tabId) {
        console.log("Switching to tab:", tabId);

        const tabs = document.querySelectorAll(".tab-content");
        const buttons = document.querySelectorAll(".tab-btn");

        tabs.forEach((tab) => tab.classList.remove("active"));
        buttons.forEach((btn) => btn.classList.remove("active"));

        document.getElementById(tabId).classList.add("active");
        document
            .querySelector(`.tab-btn[onclick="showTab('${tabId}')"]`)
            .classList.add("active");
    };

    // Open Add Credential Modal
    window.openAddCredentialModal = function () {
        console.log("Opening Add Credential Modal");
        document.getElementById("addCredentialModal").style.display = "flex";
    };

    // Close Add Credential Modal
    window.closeAddCredentialModal = function () {
        console.log("Closing Add Credential Modal");
        document.getElementById("addCredentialModal").style.display = "none";
    };

    // Open PDF Modal
    window.showModal = function (fileData) {
        console.log("Showing PDF Modal");
        const modal = document.getElementById("pdfModal");
        const viewer = document.getElementById("pdfViewer");
        viewer.src = `data:application/pdf;base64,${fileData}`;
        modal.style.display = "flex";
    };

    // Close PDF Modal
    window.closeModal = function () {
        console.log("Closing PDF Modal");
        document.getElementById("pdfModal").style.display = "none";
    };

    // Close modals on outside click
    document.addEventListener("click", function (e) {
        if (e.target.classList.contains("modal") || e.target.classList.contains("addCreModal")) {
            closeAddCredentialModal();
            closeModal();
        }
    });
});
