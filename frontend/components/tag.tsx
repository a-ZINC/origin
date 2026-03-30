const map: Record<string, string> = {
  go:         "bg-[rgba(6,182,212,0.12)]  text-[#67e8f9]",
  java:       "bg-[rgba(245,158,11,0.12)] text-[#fcd34d]",
  kubernetes: "bg-[rgba(96,165,250,0.12)] text-[#93c5fd]",
  k8s:        "bg-[rgba(96,165,250,0.12)] text-[#93c5fd]",
  kafka:      "bg-[rgba(16,185,129,0.12)] text-[#6ee7b7]",
  typescript: "bg-[rgba(49,120,198,0.12)] text-[#93c5fd]",
  docker:     "bg-[rgba(36,150,237,0.12)] text-[#7dd3fc]",
  default:    "bg-[rgba(139,92,246,0.10)] text-[#c4b5fd]",
};

export default function Tag({ name }: { name: string }) {
  const cls = map[name.toLowerCase()] ?? map.default;
  return (
    <span className={`text-[11px] px-2.5 py-0.5 rounded-[6px]
                      font-mono tracking-[0.4px] ${cls}`}>
      {name}
    </span>
  );
}