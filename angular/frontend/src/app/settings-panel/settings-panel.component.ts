import { Component, EventEmitter, Output } from '@angular/core';
import { SearchbarComponent } from '../searchbar/searchbar.component';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { FontCardComponent } from '../font-card/font-card.component';

export const previewTexts = {
  lorem:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
  alphabet: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz',
  numbers: '0123456789'
};

@Component({
  selector: 'app-settings-panel',
  standalone: true,
  imports: [SearchbarComponent, ReactiveFormsModule, FontCardComponent],
  templateUrl: './settings-panel.component.html',
  styleUrl: './settings-panel.component.scss'
})
export class SettingsPanelComponent {
  previewSize = new FormControl(24);
  previewText = new FormControl('', { updateOn: 'change' });
  sansCheckbox = new FormControl(true);
  serifCheckbox = new FormControl(true);
  handwritingCheckbox = new FormControl(true);
  displayCheckbox = new FormControl(true);
  monospaceCheckbox = new FormControl(true);

  previewTextType = 'lorem';

  @Output() previewSizeChange = new EventEmitter<number>();
  @Output() previewTextChange = new EventEmitter<string>();
  @Output() sansCheckboxChange = new EventEmitter<boolean>();
  @Output() serifCheckboxChange = new EventEmitter<boolean>();
  @Output() handwritingCheckboxChange = new EventEmitter<boolean>();
  @Output() displayCheckboxChange = new EventEmitter<boolean>();
  @Output() monospaceCheckboxChange = new EventEmitter<boolean>();
  @Output() search = new EventEmitter<string>();

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
