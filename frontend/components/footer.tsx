const links = [
  { label: "github",   href: "https://github.com/a-ZINC" },
  { label: "linkedin", href: "https://linkedin.com/in/ajinkya-singh" },
  { label: "dev.to",   href: "https://dev.to/azinc" },
];

export default function Footer() {
  return (
    <footer className="border-t border-b1 py-8 mt-8
                       flex items-center justify-between">
      <div>
        <div className="text-[15px] font-medium text-t1 mb-1">Ajinkya Singh</div>
        <div className="text-[13px] text-t3 font-mono">
          backend engineer · bengaluru, india
        </div>
      </div>

      <div className="flex items-center gap-2">
        {links.map((l) => (
          <a
            key={l.label}
            href={l.href}
            target="_blank"
            rel="noopener noreferrer"
            className="px-4 py-2 rounded-lg border border-b2 text-[13px] text-t2
                       hover:text-t1 hover:border-b3 transition-all"
          >
            {l.label}
          </a>
        ))}
        <a
          href="mailto:azinc28@gmail.com"
          className="px-4 py-2 rounded-lg bg-[var(--p)] text-[13px] text-white
                     hover:bg-[var(--p2)] transition-colors"
        >
          email me
        </a>
      </div>
    </footer>
  );
}