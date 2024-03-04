import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { ApiService } from '../../api/services/api.service';
import { Router } from '@angular/router';
import { Webfont } from '../../api/models/webfont';

@Component({
  selector: 'app-edit-font-dialog',
  templateUrl: './edit-font-dialog.component.html',
  styleUrl: './edit-font-dialog.component.scss'
})
export class EditFontDialogComponent {
  @Input() font!: Webfont;

  @Output() saved = new EventEmitter();

  editFontForm = new FormGroup({
    license: new FormControl('', { nonNullable: false }),
    category: new FormControl('Sans Serif', { nonNullable: true }),
    description: new FormControl('', { nonNullable: false })
  });

  isOpen = false;

  constructor(
    private apiClient: ApiService,
    private router: Router
  ) {}

  open() {
    this.isOpen = true;
    this.editFontForm.get('license')?.setValue(this.font.license);
    this.editFontForm.get('category')?.setValue(this.font.category);
    this.editFontForm.get('description')?.setValue(this.font.description ?? null);
  }

  updateFont() {
    if (this.editFontForm.invalid) {
      return;
    }

    const license = this.editFontForm.get('license')?.value ?? '';
    const category = this.editFontForm.get('category')?.value ?? '';
    const description = this.editFontForm.get('description')?.value ?? '';

    this.apiClient
      .updateFontByName({
        fontName: this.font.name,
        body: { license, description, category }
      })
      .subscribe(() => {
        this.isOpen = false;
        this.saved.emit();
      });
  }
}
