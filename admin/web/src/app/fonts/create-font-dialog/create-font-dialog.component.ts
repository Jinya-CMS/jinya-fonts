import { Component, Input } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { ApiService } from '../../api/services/api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-create-font-dialog',
  templateUrl: './create-font-dialog.component.html',
  styleUrl: './create-font-dialog.component.scss'
})
export class CreateFontDialogComponent {
  createFontForm = new FormGroup({
    name: new FormControl('', { nonNullable: true }),
    license: new FormControl('', { nonNullable: false }),
    category: new FormControl('Sans Serif', { nonNullable: true }),
    description: new FormControl('', { nonNullable: false })
  });

  @Input() open = false;

  constructor(
    private apiService: ApiService,
    private router: Router
  ) {}

  createFont() {
    if (this.createFontForm.invalid) {
      return;
    }

    const name = this.createFontForm.get('name')?.value ?? '';
    const license = this.createFontForm.get('license')?.value ?? '';
    const category = this.createFontForm.get('category')?.value ?? '';
    const description = this.createFontForm.get('description')?.value ?? '';

    this.apiService.createNewFont({ body: { name, license, description, category } }).subscribe((value) => {
      this.router.navigate(['/font/detail', value.name]);
    });
  }
}
