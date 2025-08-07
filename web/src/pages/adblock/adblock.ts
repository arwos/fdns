import {Component, inject, OnInit} from '@angular/core';
import {RequestService} from '@onega-ui/core';
import {AdBlockList} from './model';
import {NgClass} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {Search1Pipe} from './search1-pipe';
import {MarkPipe} from './mark-pipe';
import {TabModule} from '@onega-ui/kit';

@Component({
  selector: 'app-adblock',
  imports: [
    NgClass,
    FormsModule,
    Search1Pipe,
    MarkPipe,
    TabModule
  ],
  templateUrl: './adblock.html',
  styleUrl: './adblock.scss'
})
export class Adblock implements OnInit {
  #req = inject(RequestService)

  list: AdBlockList = {links: [], rules: []}
  search = ''

  ngOnInit() {
    this.#req.get<AdBlockList>('/adblock/list').subscribe(value => {
      this.list = value
    })
  }
}
