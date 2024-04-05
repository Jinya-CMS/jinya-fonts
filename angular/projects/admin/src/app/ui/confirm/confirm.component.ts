import { booleanAttribute, Component, EventEmitter, Input, Output } from '@angular/core';
import { EditorComponent } from '@tinymce/tinymce-angular';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-confirm',
  standalone: true,
  imports: [EditorComponent, FormsModule, ReactiveFormsModule],
  templateUrl: './confirm.component.html'
})
export class ConfirmComponent {
  @Input() open!: boolean;
  @Input({ transform: booleanAttribute }) negative!: boolean;
  @Input() title!: string;
  @Input() confirmLabel!: string;
  @Input() declineLabel!: string;

  @Output() confirm = new EventEmitter();
  @Output() decline = new EventEmitter();
}
