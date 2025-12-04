import { useMutation } from "@tanstack/react-query";
import { createMessage } from "../../api/message";

export function useCreateMessage() {
    return useMutation({
        mutationFn: createMessage,
    })
}