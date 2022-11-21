import { Component, Input, OnInit } from '@angular/core';
import { BountyHunter } from 'src/app/models/Empire';

@Component({
  selector: 'app-bounty-hunter-info',
  templateUrl: './bounty-hunter-info.component.html',
  styleUrls: ['./bounty-hunter-info.component.sass']
})
export class BountyHunterInfoComponent implements OnInit {
  @Input() info: BountyHunter

  constructor() { }

  ngOnInit(): void {
  }

}
