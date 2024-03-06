import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-searchbar',
  templateUrl: './searchbar.component.html',
  styleUrl: './searchbar.component.scss'
})
export class SearchbarComponent {
  @Output() searchTextChanged = new EventEmitter<string>();

  search($event: Event) {
    const input = $event.target as HTMLInputElement;
    this.searchTextChanged.emit(input.value);
  }
}
