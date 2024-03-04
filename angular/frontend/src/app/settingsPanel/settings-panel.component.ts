import { Component, EventEmitter, Output } from '@angular/core';
import { FontCard } from '../fontCard';
import { SearchbarComponent } from '../searchbar/searchbar.component';
import { FormControl, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-settings-panel',
  standalone: true,
  imports: [SearchbarComponent, ReactiveFormsModule],
  templateUrl: './settings-panel.component.html',
  styleUrl: './settings-panel.component.scss'
})
export class SettingsPanelComponent {
  fontList: FontCard[] = [];
  filteredFontList: FontCard[] = [];
  previewSize = new FormControl(24);
  previewText = new FormControl('');
  previewTextType = new FormControl('lorem');
  sansCheckbox = new FormControl(true);
  serifCheckbox = new FormControl(true);
  handwritingCheckbox = new FormControl(true);
  displayCheckbox = new FormControl(true);
  monospaceCheckbox = new FormControl(true);

  @Output() previewSizeChange = new EventEmitter<number>();
  @Output() previewTextChange = new EventEmitter<string>();
  @Output() sansCheckboxChange = new EventEmitter<boolean>();
  @Output() serifCheckboxChange = new EventEmitter<boolean>();
  @Output() handwritingCheckboxChange = new EventEmitter<boolean>();
  @Output() displayCheckboxChange = new EventEmitter<boolean>();
  @Output() monospaceCheckboxChange = new EventEmitter<boolean>();

  updateText() {
    switch (this.previewTextType.value) {
      case 'alphabet':
        this.previewTextChange.emit('ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz');
        break;
      case 'lorem':
        this.previewTextChange.emit(
          'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.'
        );
        break;
      case 'numbers':
        this.previewTextChange.emit('1234567890');
        break;
      default:
        this.previewTextChange.emit(this.previewText.value ?? '');
        break;
    }
  }
}
