"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/",        label: "home"     },
  { href: "/blog",    label: "blog"     },
  { href: "/series",  label: "series"   },
  { href: "/projects",label: "projects" },
  { href: "/about",   label: "about"    },
];

export default function Nav() {
  const path = usePathname();

  return (
    <nav className="flex items-center justify-between py-6 border-b border-b1
                    sticky top-0 z-50 bg-[rgba(9,9,15,0.85)] backdrop-blur-xl -mx-8 px-8">
      <Link href="/" className="flex items-center gap-2.5 no-underline">
        <div className="w-9 h-9 rounded-[10px] flex items-center justify-center
                        font-mono text-[13px] font-medium text-white
                        bg-gradient-to-br from-[#8b5cf6] to-[#06b6d4]">
          AZ
        </div>
        <span className="text-[15px] font-medium text-t1 tracking-tight">
          a-zinc<span className="text-p2">.dev</span>
        </span>
      </Link>

      <div className="flex items-center gap-1">
        {links.map((l) => {
          const active = l.href === "/"
            ? path === "/"
            : path.startsWith(l.href);
          return (
            <Link
              key={l.href}
              href={l.href}
              className={`text-[14px] px-3.5 py-1.5 rounded-lg transition-all
                ${active
                  ? "text-t1 bg-s2"
                  : "text-t3 hover:text-t2 hover:bg-s1"
                }`}
            >
              {l.label}
            </Link>
          );
        })}
      </div>

      <a
        href="mailto:azinc28@gmail.com"
        className="text-[13px] px-4 py-2 rounded-lg bg-[var(--p)] text-white
                   hover:bg-[var(--p2)] transition-colors font-medium"
      >
        hire me
      </a>
    </nav>
  );
}