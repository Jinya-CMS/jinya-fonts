import { Component } from '@angular/core';
import {RouterLink, RouterOutlet} from '@angular/router';
import {SearchbarComponent} from "./searchbar/searchbar.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'jinya-fonts';
}
