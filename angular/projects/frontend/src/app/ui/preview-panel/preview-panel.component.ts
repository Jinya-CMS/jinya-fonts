import { Component, EventEmitter, input, Output } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';

export const previewTexts = {
  lorem:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
  alphabet: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz',
  numbers: '0123456789'
};

@Component({
  selector: 'app-preview-panel',
  imports: [ReactiveFormsModule],
  templateUrl: './preview-panel.component.html',
  styleUrl: './preview-panel.component.scss'
})
export class PreviewPanelComponent {
  @Output() previewSizeChange = new EventEmitter<number>();
  @Output() previewTextChange = new EventEmitter<string>();

  previewSize = new FormControl(24);
  previewText = new FormControl(previewTexts.lorem);

  previewTextType = 'lorem';

  protected readonly input = input;

  constructor() {
    this.previewText.registerOnChange((value: string) => this.previewTextChange.emit(value));
    this.setPreviewText('lorem');
  }

  updateText(event: Event) {
    const select = event.target as HTMLSelectElement;
    // @ts-expect-error This is valid
    this.setPreviewText(select.value);
  }

  updatePreviewText(event: Event) {
    const input = event.target as HTMLInputElement;
    this.previewText.setValue(input.value);
    this.setPreviewText('custom');
  }

  setPreviewText(textType: 'custom' | 'lorem' | 'alphabet' | 'numbers') {
    this.previewTextType = textType;
    if (textType != 'custom') {
      this.previewText.setValue(previewTexts[textType]);
    }
  }
}
