---
type: playbook
status: stable
last_reviewed: 2026-05-08
sources: []
---

# Quando rebalancear exposição BRL/USD

## TL;DR
Rebalancear por **bandas** (desvio relativo da meta), não por calendário.
Para uma carteira pequena, banda de ±5 pontos percentuais costuma minimizar
custo tributário sem deixar a alocação derivar.

## Passo a passo
1. Defina a meta de exposição USD em função do horizonte e do conforto com
   volatilidade — ver [[concepts/exposicao-cambial]]. Heurística:
   15-30% para investidores cuja despesa é 100% em BRL.
2. Calcule a exposição efetiva atual (não só os ativos USD diretos —
   inclua exportadoras, se relevante).
3. Compare com a meta. Se o desvio for menor que ±5 p.p., não faça nada.
4. Se o desvio for maior, rebalanceie aportando no lado faltante antes de
   vender o lado em excesso — vendas geram fato gerador de IR.
5. Use [[entities/B3]] (IVVB11, BDRs) quando o objetivo for exposição
   simples; conta no exterior quando precisar de ETFs específicos ou já
   tenha relevância tributária americana ([[concepts/PFIC]]).

## Quando NÃO seguir esse playbook
- Volatilidade cambial extrema (movimento >5% em poucos dias): aguardar
  estabilização antes de rebalancear costuma reduzir slippage.
- Mudança de regime tributário: pausar até entender o impacto.

## Relacionado
- [[concepts/exposicao-cambial]]
- [[entities/USD]]
- [[entities/BRL]]
