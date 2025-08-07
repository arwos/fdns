export class AdBlockList {
  links!: AdBlockLinkModel[]
  rules!: AdBlockRuleModel[]
}

export class AdBlockLinkModel {
  id!: number
  link!: string
  disabled!: boolean
  updatedAt!: string
}

export class AdBlockRuleModel {
  id!: number
  linkId!: number
  rule!: string
  disabled!: boolean
  updatedAt!: string
}
