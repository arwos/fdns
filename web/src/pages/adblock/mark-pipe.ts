import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'mark'
})
export class MarkPipe implements PipeTransform {

  transform(value: string, txt: string): string {
    return value.replace(txt, `<mark>${txt}</mark>`);
  }

}
