import {Pipe, PipeTransform} from '@angular/core';
import {AdBlockRuleModel} from './model';

@Pipe({
  name: 'search1'
})
export class Search1Pipe implements PipeTransform {

  transform(value: AdBlockRuleModel[], linkId: number, search: string): AdBlockRuleModel[] {
    if (search.length < 3) {
      return [];
    }

    return value.filter(value1 => {
      return value1.linkId === linkId && value1.rule.includes(search)
    });
  }

}
