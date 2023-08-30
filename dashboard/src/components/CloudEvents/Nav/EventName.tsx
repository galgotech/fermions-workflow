import { useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { BsFillFileCodeFill } from "react-icons/bs";

import { ButtonDropdown } from "../../ButtonDropdown/ButtonDropdown";
import { CloudEventModel } from "../states/CloudEventModel";


interface Props {
    cloudEvent: CloudEventModel;
    open: (value: CloudEventModel) => void;
    save: (value: CloudEventModel) => void;
}

interface UpdateNameDTO {
  name: string;
}

export const EventName = ({ cloudEvent, save, open }: Props) => {
  const [ editName, setEditName ] = useState(false);

  const { register, handleSubmit } = useForm({ mode: 'onChange', defaultValues: { name: cloudEvent?.name } });
  const formRef = useRef<HTMLFormElement>(null)

  const onSubmit = (data: UpdateNameDTO) => {
    if (!editName) {
      return;
    }
    cloudEvent.name = data.name;
    save(cloudEvent);
    setEditName(false);
  };

  return (
    <div className="pb-2">
      <div className="d-flex justify-content-between">
        <div>
          <BsFillFileCodeFill />
          <span
            className="ps-2 fw-light align-middle"
            role="button"
            onClick={() => {
              if (!editName) {
                open(cloudEvent);
              }
            }}
          >
            {editName && (
              <form
                ref={formRef}
                className="d-inline"
                onSubmit={handleSubmit(onSubmit)}
              >
                <input
                  type="text"
                  {...register("name")}
                  style={{ width: "calc(100% - 35px)" }}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                        formRef.current?.btSubmit?.click();
                    }
                  }}
                />
                <input type="submit" name="btSubmit" className="d-none" />
              </form>
            )}
            {!editName && cloudEvent.name}
          </span>
        </div>
        <span className="nav-item dropdown" >
          <ButtonDropdown
            rename={() => setEditName(true)}
            del={() => {console.log(cloudEvent)}}
          />
        </span>
      </div>
    </div>
  );
};
