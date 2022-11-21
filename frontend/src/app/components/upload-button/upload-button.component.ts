import { Component, OnInit, EventEmitter, Output, ViewChild, ElementRef } from '@angular/core';
import { ParseEmpire } from "../../models/Empire"

@Component({
  selector: 'app-upload-button',
  templateUrl: './upload-button.component.html',
  styleUrls: ['./upload-button.component.sass']
})
export class UploadButtonComponent implements OnInit {
  @Output() newEmpire = new EventEmitter<object>();
  @Output() error = new EventEmitter();

  @ViewChild('fileInput') fileInput: ElementRef;

  constructor() { }

  ngOnInit(): void {
  }

  onFileSelected(e: any) {
    console.log("File selected : ", e)
    let file = e.target.files[0];
    if (file=="") {
      return
    }

    this.fileInput.nativeElement.value = '';

    let fileReader = new FileReader();
    fileReader.onload = (e) => {
      let read = fileReader.result;
      let content = <string>read
    
      try {
        let empire = ParseEmpire(content)
        this.newEmpire.emit(empire)
      } catch(e) {
        let errorMessage = "Could not parse this file : " + e.toString()
        this.error.emit(errorMessage)
      }
    }
    fileReader.readAsText(file);
  }

}
