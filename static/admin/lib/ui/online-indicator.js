import { localize } from '../../../lib/jinya-alpine-tools.js';

class OnlineIndicatorElement extends HTMLElement {
  constructor() {
    super();

    this.root = this.attachShadow({ mode: 'closed' });
  }

  static get observedAttributes() {
    return ['online'];
  }

  get online() {
    return this.hasAttribute('online');
  }

  set online(value) {
    if (value) {
      this.setAttribute('online', value);
    } else {
      this.removeAttribute('online');
    }
    this.#setIndicator();
  }

  connectedCallback() {
    this.root.innerHTML = `
      <style>
        .center-vertical {
          display: flex;
          align-items: center;
          gap: 0.25rem;
        }
      </style>
      <div class="center-vertical"></div>
    `;
    this.#setIndicator();
  }

  #setIndicator() {
    const element = this.root.querySelector('.center-vertical');
    if (this.online) {
      element.innerHTML = `
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="var(--positive-color)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21.801 10A10 10 0 1 1 17 3.335"/>
          <path d="m9 11 3 3L22 4"/>
        </svg>
        ${localize({ key: 'online' })}
      `;
    } else {
      element.innerHTML = `
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="var(--negative-color)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"/>
          <path d="m15 9-6 6"/>
          <path d="m9 9 6 6"/>
        </svg>
        ${localize({ key: 'offline' })}
      `;
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

if (!customElements.get('jinya-online-indicator')) {
  customElements.define('jinya-online-indicator', OnlineIndicatorElement);
}
