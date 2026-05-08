from app.services.llm_client import ChatMessage, CompletionResult, LLMClient
from app.services.wiki_reader import WikiReader

SYSTEM_PROMPT_PT = """\
Você é o assistente de investimento internacional da plataforma CBPI
(Cross-Border Portfolio Intelligence). Seu público é um investidor brasileiro
que mantém posições em BRL e USD e quer entender exposição cambial,
diversificação e implicações tributárias (IRRF, PFIC) ao investir lá fora.

Princípios:
1. Responda em português do Brasil, claro e direto.
2. Baseie-se EXCLUSIVAMENTE no conteúdo da wiki abaixo. Quando faltar
   informação na wiki, diga isso explicitamente em vez de inventar.
3. Cite as páginas da wiki que você usou no formato [[categoria/slug]] no fim
   da resposta.
4. Sempre encerre com a linha de aviso:
   "Este conteúdo é educacional e não constitui recomendação de investimento."
"""


class ChatService:
    """Wiki-grounded chat. The wiki index + schema are sent as a cached system
    block; concrete page bodies are pulled in when the user message references
    them by slug. Index-first retrieval, no embeddings.
    """

    def __init__(self, *, llm: LLMClient, wiki: WikiReader) -> None:
        self._llm = llm
        self._wiki = wiki

    async def answer(self, messages: list[ChatMessage]) -> tuple[CompletionResult, list[str]]:
        cited = self._select_pages(messages[-1].content if messages else "")
        system = self._build_system(cited)
        result = await self._llm.complete(messages, system=system)
        return result, cited

    def _build_system(self, cited: list[str]) -> str:
        parts = [SYSTEM_PROMPT_PT, "\n\n# Wiki schema\n", self._wiki.read_schema()]
        parts.append("\n\n# Wiki index\n")
        parts.append(self._wiki.read_index() or "(empty)")
        if cited:
            parts.append("\n\n# Wiki pages relevant to the current question\n")
            for rel in cited:
                page = self._wiki.read_page(rel)
                if page is None:
                    continue
                parts.append(f"\n## [[{rel.removesuffix('.md')}]]\n{page.body}\n")
        return "".join(parts)

    def _select_pages(self, query: str) -> list[str]:
        if not query:
            return []
        q = query.lower()
        hits: list[str] = []
        for rel in self._wiki.list_pages():
            slug = rel.removesuffix(".md").split("/", 1)[1].lower()
            if slug in q or slug.replace("-", " ") in q:
                hits.append(rel)
        return hits[:6]
