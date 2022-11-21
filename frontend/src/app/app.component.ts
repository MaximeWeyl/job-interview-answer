import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar'

import * as z from 'zod';
import {environment} from '../environments/environment'
import {
  trigger,
  state,
  style,
  animate,
  transition,
  keyframes
} from '@angular/animations';

import { BountyHunter, Empire, schemaEmpire } from './models/Empire';

const schemaNode = z.object({
  "FuelLeft": z.number().int(),
  "PlanetName": z.string(),
  "Time": z.number().int()
})

const schemaResponse = z.object({
  odds: z.number(),
  path: z.array(schemaNode)
})


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass'],
  animations: [
    trigger('oddsAnimation', [
      state('turn1', style({
        transform: 'rotateZ(0deg)',
      })),
      state('turn2', style({
        transform: 'rotateZ(360deg)',
      })),
      transition('turn1 => turn2', [
        animate('0.5s ease-in-out')
      ]),
      transition('turn2 => turn1', [
        animate('0.5s ease-in-out')
      ]),
    ])
  ]
})
export class AppComponent {
  title = 'frontend';
  odds = 0 ;
  oddsAnimationState = "turn1"

  empireInfo: Empire
  bounty_hunter: BountyHunter // TODO : remove this one 

  constructor(private http: HttpClient, private snackBar: MatSnackBar) {
  }

  onNewEmpire(empireInfo: object) {
    let request = {
      empire: empireInfo
    }

    // Make the API call
    this.http.post<any>(environment.api_url+"/give-me-the-odds", request).subscribe(data => {
      console.log(data)

      try {
        schemaResponse.parse(data)
        this.snackBar.open("Success" + data.toString(), 'End now', {
          duration: 500,
        });

        this.empireInfo = schemaEmpire.parse(empireInfo) // Throws if not valid
        
        try {
          this.bounty_hunter = this.empireInfo.bounty_hunters[0]
        } catch {
          this.bounty_hunter = new BountyHunter()
        }

        this.odds = data.odds * 100

        if (this.oddsAnimationState == "turn1"){
          this.oddsAnimationState = "turn2"
        } else {
          this.oddsAnimationState = "turn1"
        }


      } catch(e) {
        this.onError("Could not parse data returned by the API : " + e.toString())
      }
    })
  }

  onError(errorMessage: string) {
    this.snackBar.open(errorMessage, 'End now', {
      duration: 10000,
    });
  }
}
