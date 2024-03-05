import { Component, EventEmitter, input, Output } from '@angular/core';
import { SearchbarComponent } from '../searchbar/searchbar.component';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { FontCardComponent } from '../font-card/font-card.component';

export const previewTexts = {
  lorem:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
  alphabet: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz',
  numbers: '0123456789'
};

export interface WebfontFilter {
  searchText: string;
  sansSerif: boolean;
  serif: boolean;
  handwriting: boolean;
  display: boolean;
  monospace: boolean;
}

@Component({
  selector: 'app-settings-panel',
  standalone: true,
  imports: [SearchbarComponent, ReactiveFormsModule, FontCardComponent],
  templateUrl: './settings-panel.component.html',
  styleUrl: './settings-panel.component.scss'
})
export class SettingsPanelComponent {
  @Output() previewSizeChange = new EventEmitter<number>();
  @Output() previewTextChange = new EventEmitter<string>();
  @Output() filter = new EventEmitter<WebfontFilter>();

  previewSize = new FormControl(24);
  previewText = new FormControl('');
  sansCheckbox = new FormControl(true);
  serifCheckbox = new FormControl(true);
  handwritingCheckbox = new FormControl(true);
  displayCheckbox = new FormControl(true);
  monospaceCheckbox = new FormControl(true);

  previewTextType = 'lorem';

  activeFilter: WebfontFilter = {
    sansSerif: true,
    serif: true,
    handwriting: true,
    display: true,
    monospace: true,
    searchText: ''
  };

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

  allCheckboxesChecked() {
    if (
      !this.activeFilter.sansSerif &&
      !this.activeFilter.serif &&
      !this.activeFilter.handwriting &&
      !this.activeFilter.display &&
      !this.activeFilter.monospace
    ) {
      this.sansCheckbox.setValue(true);
      this.serifCheckbox.setValue(true);
      this.handwritingCheckbox.setValue(true);
      this.displayCheckbox.setValue(true);
      this.monospaceCheckbox.setValue(true);

      this.activeFilter.sansSerif = true;
      this.activeFilter.serif = true;
      this.activeFilter.handwriting = true;
      this.activeFilter.display = true;
      this.activeFilter.monospace = true;

      this.filter.emit(this.activeFilter);
    }
  }

  sansCheckboxChanged($event: Event) {
    const input = $event.target as HTMLInputElement;
    this.activeFilter.sansSerif = input.checked;
    this.filter.emit(this.activeFilter);
    this.allCheckboxesChecked();
  }

  serifCheckboxChanged($event: Event) {
    const input = $event.target as HTMLInputElement;
    this.activeFilter.serif = input.checked;
    this.filter.emit(this.activeFilter);
    this.allCheckboxesChecked();
  }

  handwritingCheckboxChanged($event: Event) {
    const input = $event.target as HTMLInputElement;
    this.activeFilter.handwriting = input.checked;
    this.filter.emit(this.activeFilter);
    this.allCheckboxesChecked();
  }

  displayCheckboxChanged($event: Event) {
    const input = $event.target as HTMLInputElement;
    this.activeFilter.display = input.checked;
    this.filter.emit(this.activeFilter);
    this.allCheckboxesChecked();
  }

  monospaceCheckboxChanged($event: Event) {
    const input = $event.target as HTMLInputElement;
    this.activeFilter.monospace = input.checked;
    this.filter.emit(this.activeFilter);
    this.allCheckboxesChecked();
  }

  search(term: string) {
    this.activeFilter.searchText = term;
    this.filter.emit(this.activeFilter);
  }
}
