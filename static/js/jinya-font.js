function updateBody() {
    const previewSize = document.getElementById('previewsize').value;
    const previewText = document.getElementById('previewtext').value;
    const previewType = document.getElementById('previewtype').value;
    document.querySelectorAll('[data-role=preview]').forEach((item) => {
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
    });
}

function updateLinks() {
    const urlsearch = new URLSearchParams(location.search);
    const basepath = `${location.origin}/css2?family=${urlsearch.get('font')}`;

    const selectedStyles = [];
    document.querySelectorAll('[data-role=select]:checked').forEach((item) => {
        selectedStyles.push({
            style: item.getAttribute('data-style') === 'italic' ? 1 : 0,
            weight: item.getAttribute('data-weight'),
        });
    });

    let url;
    if (selectedStyles.length === 0) {
        url = basepath;
    } else {
        url = `${basepath}:ital,wght@${selectedStyles.map(item => `${item.style},${item.weight}`).join('%3B')}`;
    }

    document.getElementById('font-html-link').innerText = `<link rel="stylesheet" type="text/css" href="${url}">`;
    document.getElementById('font-download-link').setAttribute('href', `/download?font=${urlsearch.get('font')}`);
}

function dedupe(arr) {
    const hashTable = {};

    return arr.filter(function (el) {
        const key = JSON.stringify(el);
        const match = Boolean(hashTable[key]);

        return (match ? false : hashTable[key] = true);
    });
}

document.addEventListener('DOMContentLoaded', async () => {
        const urlsearch = new URLSearchParams(location.search);
        const response = await fetch(`/api/font?font=${urlsearch.get('font')}`);
        const font = await response.json();
        const filteredFontStylesAndWeights = dedupe(font.fonts.map(m => ({
            weight: m.weight,
            style: m.style,
        }))).sort((a, b) => {
            const weighta = a.weight;
            const weightb = b.weight;
            const stylea = a.style;

            if (weighta < weightb) {
                return -1;
            }

            if (weighta > weightb) {
                return 1;
            }

            if (stylea.toLowerCase() === 'italic') {
                return 1;
            }

            return 0;
        });

        const cssFont = filteredFontStylesAndWeights.map(m => `${m.style.toLowerCase() === 'italic' ? '1' : '0'},${m.weight}`).join('%3B');
        document.head.insertAdjacentHTML('beforeend', `<link type="text/css" rel="stylesheet" href="/css2?family=${encodeURI(font.name)}:ital,wght@${cssFont}&display=block">`);

        const contentTmpl = document.getElementById('content-tmpl').innerHTML;
        const stylePreviewTmpl = document.getElementById('style-preview').innerHTML;
        const designerDetailsTmpl = document.getElementById('designer-details').innerHTML;
        const designerNames = font.designers?.map(designer => designer.name).join(', ') ?? '';
        const content = document.getElementById('content');
        let styleHtml = '';

        for (const styleAndWeight of filteredFontStylesAndWeights) {
            styleHtml += stylePreviewTmpl
                .replaceAll('#= title #', `${font.name} ${styleAndWeight.weight} ${styleAndWeight.style}`)
                .replaceAll('#= weight #', styleAndWeight.weight)
                .replaceAll('#= style #', styleAndWeight.style)
                .replaceAll('#= family #', font.name);
        }

        let designersHtml = '';
        for (const designer of font.designers ?? []) {
            designersHtml += designerDetailsTmpl
                .replaceAll('#= designer #', designer.name)
                .replaceAll('#= bio #', designer.bio);
        }

        content.innerHTML = contentTmpl
            .replaceAll('#= name #', font.name)
            .replaceAll('#= designers #', designerNames)
            .replaceAll('#= about #', font.description.replaceAll('<a', '<a target="_blank"'))
            .replaceAll('#= style-preview-container #', styleHtml)
            .replaceAll('#= style-designers #', designersHtml.replaceAll('<a', '<a target="_blank"'));

        document.getElementById('previewtype').addEventListener('change', updateBody);
        document.getElementById('previewtext').addEventListener('input', updateBody);
        const previewSizeText = document.getElementById('previewsizetext');
        document.getElementById('previewsize').addEventListener('input', (e) => {
            const value = e.currentTarget.value;
            updateBody();
            previewSizeText.innerText = value + 'px';
        });

        document.getElementById('font-html-link').innerText = `<link rel="stylesheet" type="text/css" href="${location.origin}/css2?family=${urlsearch.get('font')}">`;
        document.getElementById('font-css').innerText = `body {
    font-family: ${urlsearch.get('font')}, ${font.fonts[0].category};
}`;

        let license = font.license;
        switch (license) {
            case "apache2":
                license = '<a target="_blank" href="https://www.apache.org/licenses/LICENSE-2.0">Apache License, Version 2.0</a>';
                break;
            case "ofl":
                license = '<a target="_blank" href="https://openfontlicense.org/">Open Font License</a>';
                break;
            case "ufl":
                license = '<a target="_blank" href="https://font.ubuntu.com/ufl/">Ubuntu Font License</a>';
                break;
        }
        document.getElementById('font-license').innerHTML = license;

        document.querySelectorAll('[data-role=select]').forEach(item => item.addEventListener('change', updateLinks));

        updateLinks();
        updateBody();
    }
);