let idx = 0;
let originalItems = [];
let items = [];

function updateBody() {
    const previewSize = document.getElementById('previewsize').value;
    const previewText = document.getElementById('previewtext').value;
    const previewType = document.getElementById('previewtype').value;
    document.querySelectorAll('[data-role=body]').forEach((item) => {
        switch (previewType) {
            case 'custom':
                item.innerHTML = previewText;
                break;
            case 'alphabet':
                item.innerHTML = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz'
                break;
            case 'lorem':
                item.innerHTML = 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.';
                break;
            case 'numbers':
                item.innerHTML = '1234567890';
                break;
        }
        item.style.fontSize = previewSize + 'px';
        item.style.fontFamily = `'${item.getAttribute('data-font-name')}'`;
    });
}

function populate(initial) {
    if (idx + 1 === items.length) {
        return;
    }
    const template = document.getElementById('font-card').innerHTML;
    const cardlist = document.getElementById('cardlist');
    let html = '';
    let target = idx + 5;
    if (initial === true) {
        idx = 0;
        target = 48;
    }

    if (target > items.length) {
        target = items.length;
    }

    for (; idx < target; idx++) {
        const item = items[idx];
        const designers = item.designers?.map(designer => designer.name).join(', ') ?? '';
        html += template.replaceAll('#= fontname #', item.name)
            .replaceAll('#= title #', item.name)
            .replaceAll('#= subtitle #', designers);
    }
    cardlist.classList.remove('jinya-card__list--loading');
    if (initial === true) {
        cardlist.innerHTML = '';
    }
    cardlist.insertAdjacentHTML('beforeend', html);

    updateBody();
}

function filterItems(keyword) {
    const sans = document.getElementById('sans').checked;
    const serif = document.getElementById('serif').checked;
    const handwritten = document.getElementById('handwritten').checked;
    const display = document.getElementById('display').checked;
    const monospace = document.getElementById('monospace').checked;

    items = originalItems.filter((item) => {
        let result = false;
        if (sans) {
            result |= item.category.toLowerCase() === 'sans serif';
        }
        if (serif) {
            result |= item.category.toLowerCase() === 'serif';
        }
        if (handwritten) {
            result |= item.category.toLowerCase() === 'handwriting';
        }
        if (display) {
            result |= item.category.toLowerCase() === 'display';
        }
        if (monospace) {
            result |= item.category.toLowerCase() === 'monospace';
        }

        return result;
    }).filter(item => item.name.toLowerCase().includes(keyword));
}

document.addEventListener('DOMContentLoaded', async () => {
    const response = await fetch('/api/font');
    items = await response.json();
    originalItems = items;
    window.addEventListener('scroll', populate);
    let styles = '';
    for (const item of items) {
        styles += `<link type="text/css" rel="stylesheet" href="/css2?family=${encodeURI(item.name)}&display=swap">`
    }
    document.head.insertAdjacentHTML('beforeend', styles);

    document.getElementById('cardlist').innerHTML = '';

    const cardlist = document.getElementById('cardlist');
    const searchbox = document.getElementById('searchbox');
    searchbox.addEventListener('keyup', () => {
        cardlist.innerHTML = '<div class="jinya-loader"></div>';
        cardlist.classList.add('jinya-card__list--loading');

        filterItems(searchbox.value.toLowerCase());
        populate(true);
    });

    document.getElementById('previewtype').addEventListener('change', () => {
        updateBody();
    });
    document.getElementById('previewtext').addEventListener('input', () => {
        updateBody();
    });
    const previewSizeText = document.getElementById('previewsizetext');
    document.getElementById('previewsize').addEventListener('input', (e) => {
        const value = e.currentTarget.value;
        updateBody();
        previewSizeText.innerText = value + 'px';
    })

    document.getElementById('sans').addEventListener('change', () => {
        filterItems(searchbox.value.toLowerCase());
        populate(true);
    });
    document.getElementById('serif').addEventListener('change', () => {
        filterItems(searchbox.value.toLowerCase());
        populate(true);
    });
    document.getElementById('handwritten').addEventListener('change', () => {
        filterItems(searchbox.value.toLowerCase());
        populate(true);
    });
    document.getElementById('display').addEventListener('change', () => {
        filterItems(searchbox.value.toLowerCase());
        populate(true);
    });
    document.getElementById('monospace').addEventListener('change', () => {
        filterItems(searchbox.value.toLowerCase());
        populate(true);
    });

    populate(true);
    updateBody();
});