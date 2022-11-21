import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { UploadButtonComponent } from './components/upload-button/upload-button.component';
import { MatSnackBarModule } from '@angular/material/snack-bar'

import { FlexLayoutModule } from '@angular/flex-layout';
import { HttpClientModule } from '@angular/common/http';
import { BountyHunterInfoComponent } from './components/bounty-hunter-info/bounty-hunter-info.component';


@NgModule({
  declarations: [
    AppComponent,
    UploadButtonComponent,
    BountyHunterInfoComponent
  ],
  imports: [
    FlexLayoutModule,
    BrowserModule,
    BrowserAnimationsModule,
    MatIconModule,
    MatButtonModule,
    MatSnackBarModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
