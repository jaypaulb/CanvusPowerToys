/**
 * RCU Page JavaScript
 * Handles user identification, note posting, file upload, and QR code generation
 */

// Color utility function
function isColorDark(hexColor) {
    const r = parseInt(hexColor.slice(1, 3), 16);
    const g = parseInt(hexColor.slice(3, 5), 16);
    const b = parseInt(hexColor.slice(5, 7), 16);
    return (r * 0.299 + g * 0.587 + b * 0.114) / 255 < 0.5;
}

document.addEventListener('DOMContentLoaded', () => {
    // User Identification Elements
    const teamButtons = document.querySelectorAll('.team-button');
    const userNameInput = document.getElementById('username');
    const userIdentification = document.getElementById('user-identification');
    const identificationMessage = document.getElementById('message');
    const userDashboard = document.getElementById('user-dashboard');
    const welcomeMessage = document.getElementById('welcomeMessage');
    const noteSquare = document.getElementById('noteSquare');
    const uploadItemButton = document.getElementById('uploadItemButton');
    const postNoteButton = document.getElementById('postNoteButton');
    const qrCodeDiv = document.getElementById('qrcode');

    // Check if required elements exist
    if (!userNameInput || !userIdentification || !userDashboard) {
        console.error('RCU: Required elements not found');
        return;
    }

    // Variables to store user info
    let selectedTeam = null;
    let userName = null;
    let userColor = null;

    // Team colors (ROYGBIV) - exact RGB values as specified
    const teamColors = {
        1: "rgb(255, 0, 0)",      // Red
        2: "rgb(255, 127, 0)",    // Orange
        3: "rgb(255, 255, 0)",    // Yellow
        4: "rgb(0, 255, 0)",      // Green
        5: "rgb(0, 0, 255)",      // Blue
        6: "rgb(75, 0, 130)",     // Indigo
        7: "rgb(139, 0, 255)"     // Violet
    };

    // Text colors for each team (matching the provided style)
    const teamTextColors = {
        1: "rgb(255, 255, 255)",  // White for Red
        2: "rgb(0, 0, 0)",         // Black for Orange
        3: "rgb(0, 0, 0)",         // Black for Yellow
        4: "rgb(0, 0, 0)",         // Black for Green
        5: "rgb(255, 255, 255)",  // White for Blue
        6: "rgb(255, 255, 255)",  // White for Indigo
        7: "rgb(255, 255, 255)"    // White for Violet
    };

    // Initialize team button colors - match exact style provided
    teamButtons.forEach(button => {
        const teamNumber = parseInt(button.dataset.team);
        if (teamNumber && teamColors[teamNumber]) {
            button.style.backgroundColor = teamColors[teamNumber];
            button.style.color = teamTextColors[teamNumber];
        }
    });

    // Generate QR code - use actual server IP address (NOT localhost)
    async function generateQRCode() {
        if (!qrCodeDiv) return;

        let qrUrl = null;

        // MUST get server info from API to get actual IP address
        try {
            const response = await fetch('/api/server-info');
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            const data = await response.json();

            // Use the actual IP address from server
            if (data.ip && data.ip !== 'localhost' && data.ip !== '127.0.0.1') {
                const protocol = window.location.protocol;
                const port = data.port || window.location.port || '8080';
                qrUrl = `${protocol}//${data.ip}:${port}/rcu.html`;
            } else if (data.url && !data.url.includes('localhost') && !data.url.includes('127.0.0.1')) {
                qrUrl = `${data.url}/rcu.html`;
            } else if (data.hostname && data.hostname !== 'localhost' && data.hostname !== '127.0.0.1') {
                const protocol = window.location.protocol;
                const port = data.port || window.location.port || '8080';
                qrUrl = `${protocol}//${data.hostname}:${port}/rcu.html`;
            }
        } catch (err) {
            console.error('Failed to fetch server info for QR code:', err);
        }

        // If we still don't have a valid IP-based URL, show error
        if (!qrUrl || qrUrl.includes('localhost') || qrUrl.includes('127.0.0.1')) {
            qrCodeDiv.innerHTML = '<p style="color: red;">Error: Could not determine server IP address. QR code unavailable.</p>';
            return;
        }

        // Clear existing QR code
        qrCodeDiv.innerHTML = '';

        // Check if QRCode library is loaded
        if (typeof QRCode === 'undefined') {
            qrCodeDiv.innerHTML = '<p>QR Code library not loaded</p>';
            return;
        }

        // Generate QR code
        const qrcode = new QRCode(qrCodeDiv, {
            text: qrUrl,
            width: 256,
            height: 256,
            colorDark: '#000',
            colorLight: '#fff',
            correctLevel: QRCode.CorrectLevel.H
        });

        // Center the QR code
        if (qrCodeDiv.parentElement) {
            qrCodeDiv.parentElement.style.textAlign = 'center';
            qrCodeDiv.style.display = 'inline-block';
        }
    }

    // Generate QR code on page load
    generateQRCode();

    // Function to update identification message
    function updateIdentificationMessage(text, type) {
        if (identificationMessage) {
            identificationMessage.textContent = text;
            identificationMessage.className = `message ${type}`;
            identificationMessage.style.display = 'block';
        }
    }

    // Team selection and immediate identification
    teamButtons.forEach(button => {
        button.addEventListener('click', async () => {
            const name = userNameInput.value.trim();
            if (!name) {
                updateIdentificationMessage("Please enter your name first.", "error");
                return;
            }

            teamButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            selectedTeam = parseInt(button.dataset.team);
            userName = name;

            try {
                // Identify user - this endpoint should exist
                const response = await fetch('/identify-user', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ team: selectedTeam, name: userName })
                });

                const data = await response.json();
                if (response.ok && data.success) {
                    userColor = data.color;

                    // Get canvas name from API
                    try {
                        const canvasResponse = await fetch('/api/canvas/info');
                        const canvasData = await canvasResponse.json();
                        const canvasName = canvasData.canvas_name || "Unknown Canvas";

                        // Update visibility of sections
                        if (userIdentification) userIdentification.style.display = "none";
                        if (userDashboard) {
                            userDashboard.style.display = "block";
                            // Update welcome message and note square
                            if (welcomeMessage) {
                                welcomeMessage.textContent = `Welcome, ${userName}, you are currently posting to Team ${selectedTeam} on ${canvasName}.`;
                            }
                            if (noteSquare) {
                                noteSquare.style.backgroundColor = userColor;
                                noteSquare.style.color = isColorDark(userColor) ? '#FFFFFF' : '#000000';
                            }
                        }
                    } catch (err) {
                        console.error('Error fetching canvas info:', err);
                        // Continue anyway with default canvas name
                        if (userIdentification) userIdentification.style.display = "none";
                        if (userDashboard) {
                            userDashboard.style.display = "block";
                            if (welcomeMessage) {
                                welcomeMessage.textContent = `Welcome, ${userName}, you are currently posting to Team ${selectedTeam}.`;
                            }
                            if (noteSquare) {
                                noteSquare.style.backgroundColor = userColor;
                                noteSquare.style.color = isColorDark(userColor) ? '#FFFFFF' : '#000000';
                            }
                        }
                    }

                    updateIdentificationMessage("Identification successful!", "success");
                } else {
                    updateIdentificationMessage(data.error || "Identification failed.", "error");
                }
            } catch (error) {
                console.error("Error identifying user:", error);
                updateIdentificationMessage("An error occurred during identification.", "error");
            }
        });
    });

    // Handle post note button
    if (postNoteButton) {
        postNoteButton.addEventListener('click', async () => {
            if (!noteSquare) return;

            const noteText = noteSquare.textContent.trim();
            if (!noteText) {
                updateIdentificationMessage("Please enter text in the note.", "error");
                return;
            }

            if (!selectedTeam || !userName || !userColor) {
                updateIdentificationMessage("Please identify yourself first.", "error");
                return;
            }

            try {
                const response = await fetch('/create-note', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        team: selectedTeam,
                        name: userName,
                        text: noteText,
                        color: userColor
                    })
                });

                const data = await response.json();
                if (response.ok && data.success) {
                    updateIdentificationMessage("Note posted successfully!", "success");
                    noteSquare.textContent = ""; // Clear the note
                } else {
                    updateIdentificationMessage(data.error || "Failed to post note.", "error");
                }
            } catch (error) {
                console.error("Error posting note:", error);
                updateIdentificationMessage("An error occurred while posting the note.", "error");
            }
        });
    }

    // Handle file upload button
    if (uploadItemButton) {
        uploadItemButton.addEventListener('click', () => {
            if (!selectedTeam || !userName) {
                updateIdentificationMessage("Please identify yourself first.", "error");
                return;
            }

            const fileInput = document.createElement('input');
            fileInput.type = 'file';
            fileInput.accept = '.jpg,.jpeg,.png,.gif,.bmp,.tiff,.mp4,.avi,.mov,.wmv,.pdf,.mkv';

            fileInput.onchange = async () => {
                const file = fileInput.files[0];
                if (file) {
                    const formData = new FormData();
                    formData.append('team', selectedTeam);
                    formData.append('name', userName);
                    formData.append('file', file);

                    try {
                        updateIdentificationMessage("Uploading file...", "loading");
                        const response = await fetch('/upload-item', {
                            method: 'POST',
                            body: formData
                        });

                        const data = await response.json();
                        if (response.ok && data.success) {
                            updateIdentificationMessage(data.message || "File uploaded successfully!", "success");
                        } else {
                            updateIdentificationMessage(data.error || "Upload failed.", "error");
                        }
                    } catch (error) {
                        console.error("Error uploading file:", error);
                        updateIdentificationMessage("An error occurred while uploading the file.", "error");
                    }
                }
            };

            fileInput.click();
        });
    }
});
