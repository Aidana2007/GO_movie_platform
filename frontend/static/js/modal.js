// Modal utility functions for displaying messages and confirmations

// Show a message modal (info, success, error)
function showMessage(title, message, type = 'info') {
    const modal = createModal();

    // Set icon and color based on type
    let icon = '';
    let iconColor = '';
    if (type === 'success') {
        icon = '✓';
        iconColor = 'var(--success)';
    } else if (type === 'error') {
        icon = '✕';
        iconColor = 'var(--danger)';
    } else {
        icon = 'ℹ';
        iconColor = 'var(--primary)';
    }

    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <span class="modal-icon" style="color: ${iconColor};">${icon}</span>
                <h3>${title}</h3>
            </div>
            <div class="modal-body">
                <p>${message}</p>
            </div>
            <div class="modal-footer">
                <button class="btn btn-primary" onclick="closeModal()">OK</button>
            </div>
        </div>
    `;

    showModal(modal);
}

// Show a confirmation modal
function showConfirm(title, message, onConfirm, onCancel) {
    const modal = createModal();

    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <span class="modal-icon" style="color: var(--warning);">⚠</span>
                <h3>${title}</h3>
            </div>
            <div class="modal-body">
                <p>${message}</p>
            </div>
            <div class="modal-footer">
                <button class="btn btn-secondary" id="modal-cancel">Cancel</button>
                <button class="btn btn-primary" id="modal-confirm">Confirm</button>
            </div>
        </div>
    `;

    showModal(modal);

    document.getElementById('modal-confirm').addEventListener('click', () => {
        closeModal();
        if (onConfirm) onConfirm();
    });

    document.getElementById('modal-cancel').addEventListener('click', () => {
        closeModal();
        if (onCancel) onCancel();
    });
}

// Helper function to create modal element
function createModal() {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.id = 'modal-overlay';

    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.id = 'custom-modal';

    overlay.appendChild(modal);

    // Close on overlay click
    overlay.addEventListener('click', (e) => {
        if (e.target === overlay) {
            closeModal();
        }
    });

    // Close on ESC key
    document.addEventListener('keydown', handleEscKey);

    return modal;
}

// Helper function to show modal
function showModal(modal) {
    const overlay = modal.parentElement;
    document.body.appendChild(overlay);

    // Trigger animation
    setTimeout(() => {
        overlay.classList.add('active');
    }, 10);
}

// Helper function to close modal
function closeModal() {
    const overlay = document.getElementById('modal-overlay');
    if (overlay) {
        overlay.classList.remove('active');
        setTimeout(() => {
            overlay.remove();
            document.removeEventListener('keydown', handleEscKey);
        }, 300);
    }
}

// Handle ESC key press
function handleEscKey(e) {
    if (e.key === 'Escape') {
        closeModal();
    }
}
