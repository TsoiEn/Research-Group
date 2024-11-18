document.addEventListener('DOMContentLoaded', () => {
    const tableBody = document.querySelector('#student-table tbody');
    const addStudentBtn = document.getElementById('add-student-btn');

    // Add student functionality
    addStudentBtn.addEventListener('click', () => {
        const newRow = document.createElement('tr');

        const studentIdCell = document.createElement('td');
        studentIdCell.contentEditable = true;
        studentIdCell.classList.add('student-id');
        studentIdCell.style.cursor = "pointer";

        const courseCell = document.createElement('td');
        courseCell.contentEditable = true;

        const lastNameCell = document.createElement('td');
        lastNameCell.contentEditable = true;

        const firstNameCell = document.createElement('td');
        firstNameCell.contentEditable = true;

        const emailCell = document.createElement('td');
        emailCell.contentEditable = true;

        newRow.appendChild(studentIdCell);
        newRow.appendChild(courseCell);
        newRow.appendChild(lastNameCell);
        newRow.appendChild(firstNameCell);
        newRow.appendChild(emailCell);

        tableBody.appendChild(newRow);

        // Redirect on Student ID click
        studentIdCell.addEventListener('click', () => {
            const studentId = studentIdCell.textContent.trim();
            if (studentId) {
                window.location.href = `student-details.html?studentId=${studentId}`;
            }
        });
    });
});
