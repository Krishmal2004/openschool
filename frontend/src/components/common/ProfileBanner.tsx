import Avatar from "./Avatar";

interface Props {
  name: string;
  meta?: string;
  actions: React.ReactNode;
}

export default function ProfileBanner({ name, meta, actions }: Props) {
  return (
    <div className="os-profile__banner">
      <Avatar name={name} />
      <div style={{ flex: 1 }}>
        <p className="os-profile__name">{name}</p>
        {meta && <p className="os-profile__meta">{meta}</p>}
      </div>
      <div className="os-profile__actions">{actions}</div>
    </div>
  );
}
