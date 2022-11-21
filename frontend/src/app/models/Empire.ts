import * as z from 'zod';

export class BountyHunter {
    day: number
    planet: string
}

export class Empire {
    countdown: number
    bounty_hunters: BountyHunter[]
}

export const schemaBounty = z.object({
  planet: z.string(),
  day: z.number().int().positive()
})

export const schemaEmpire = z.object({
  countdown: z.number().int().positive(),
  bounty_hunters: z.array(schemaBounty)
})

export function ParseEmpire(content: string) : Empire {
    let obj = JSON.parse(content);

    obj = schemaEmpire.parse(obj) // Throws if not valid
    let empire: Empire = obj
    return empire
}