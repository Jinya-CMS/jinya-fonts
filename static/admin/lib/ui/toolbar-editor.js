import { createJodit } from './jodit.js';
import { EditorChangeEvent } from './events/EditorChangeEvent.js';

class ToolbarEditorElement extends HTMLElement {
  constructor() {
    super();

    this.root = this.attachShadow({ mode: 'closed' });
    this.editor = null;
  }

  connectedCallback() {
    this.root.innerHTML = `
      <style>
        @import "/static/lib/cosmo/typography.css";
        @import "/static/admin/lib/jodit/jodit.css";
        @import "/static/admin/css/jodit.css";
      </style>
      <textarea></textarea>
    `;
    this.editor = createJodit(this.root.querySelector('textarea'), false, false, this.height);
    this.editor.value = this.content;
    this.editor.events.on('change', (e) => {
      this.dispatchEvent(new EditorChangeEvent(e));
    });
  }

  disconnectedCallback() {
    this.editor?.destruct();
  }

  static get observedAttributes() {
    return ['content', 'focused'];
  }

  get content() {
    return this.getAttribute('content');
  }

  set content(value) {
    this.setAttribute('content', value);
    this.editor.value = value;
  }

  get height() {
    return this.getAttribute('height');
  }

  set height(value) {
    this.setAttribute('height', value);
  }

  get focused() {
    return this.hasAttribute('focused');
  }

  set focused(value) {
    if (value) {
      this.setAttribute('focused', value);
      this.editor.focus();
    } else {
      this.removeAttribute('focused');
    }
  }

  attributeChangedCallback(property, oldValue, newValue) {
    if (oldValue === newValue) {
      return;
    }

    const propertyName = property.replace(/-([a-z])/g, (m, w) => w.toUpperCase());
    this[propertyName] = newValue;
  }
}

if (!customElements.get('jinya-toolbar-editor')) {
  customElements.define('jinya-toolbar-editor', ToolbarEditorElement);
}
