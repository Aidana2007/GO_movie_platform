function showMessage(title, message, type = 'info') {
    const modal = createModal();


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

function createModal() {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.id = 'modal-overlay';

    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.id = 'custom-modal';

    overlay.appendChild(modal);

    
    overlay.addEventListener('click', (e) => {
        if (e.target === overlay) {
            closeModal();
        }
    });

    document.addEventListener('keydown', handleEscKey);

    return modal;
}

function showModal(modal) {
    const overlay = modal.parentElement;
    document.body.appendChild(overlay);

    setTimeout(() => {
        overlay.classList.add('active');
    }, 10);
}

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

function handleEscKey(e) {
    if (e.key === 'Escape') {
        closeModal();
    }
}
