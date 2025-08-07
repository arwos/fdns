import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {KitModule} from '@onega-ui/kit';
import {CoreModule} from '@onega-ui/core';


@Component({
  selector: 'app-root',
  imports: [RouterOutlet, KitModule, CoreModule],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {

}
