import { Routes } from '@angular/router';
import { StartpageComponent } from './startpage/startpage.component';

export const routes: Routes = [
  {
    component: StartpageComponent,
    path: '**'
  }
];
